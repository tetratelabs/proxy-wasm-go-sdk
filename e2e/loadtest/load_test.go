// Copyright 2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loadtest

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"testing"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/fnet"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/tetratelabs/proxy-wasm-go-sdk/e2e/testutil"
)

const (
	targetNintyninthPercentileLatencyLimit = 200 // ms
	targetSuccessRate                      = 0.7
)

var (
	qps                 = flag.Float64("qps", 0, "QPS to run load test")
	duration            = flag.Int("duration", 10, "Duration of test in seconds")
	payloadSize         = flag.Int("payloadSize", 256, "Payload size in kilo bytes")
	targetExample       = flag.String("targetExample", "http_headers", "Target example to run load test")
	memoryUsageGraphDst = flag.String("memoryUsageGraphDst", "", "Destination path for saving the memory usage graph")
)

var (
	gcStatLogFormat = regexp.MustCompile(`\[memstat\]\[contextID=(\d+)\]\[unixnanotime=(\d+)\] heap size: in-use \/ reserved = (\d+) \/ (\d+) bytes`)
)

// TestAvailabilityAgainstHighHTTPLoad tests the availability of the proxy with wasm filter against a high HTTP load
func TestAvailabilityAgainstHighHTTPLoad(t *testing.T) {
	stdErr, kill := testutil.StartEnvoyWith(*targetExample, t, 8001)
	defer kill()

	opts := fhttp.HTTPRunnerOptions{}
	opts.URL = "http://localhost:18000"
	opts.AllowInitialErrors = true
	opts.NumThreads = runtime.NumCPU()
	opts.Percentiles = []float64{99.0}
	opts.AddAndValidateExtraHeader("Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M")

	// Set payload (request body) size
	fnet.ChangeMaxPayloadSize(*payloadSize * fnet.KILOBYTE)
	opts.Payload = fnet.Payload

	fortioLog := new(bytes.Buffer)
	opts.Out = fortioLog

	opts.HTTPReqTimeOut = 5000 * time.Second // Avoid timeouts on huge payloads
	log.Printf("\tDuration = %d [s], payloadSize = %d [KB]\n", *duration, *payloadSize)
	opts.QPS = *qps * 1.5 // In order to reach the target QPS, we need to set a little bit higer target QPS.
	opts.Duration = time.Duration(*duration) * time.Second

	// Run memory profiling to find out memory stability of SDK
	err := testutil.StartEnvoyMemoryProfile(8001)
	require.NoError(t, err)

	// Start generating load
	results, err := fhttp.RunHTTPTest(&opts)
	require.NoError(t, err)

	// Stop memory profiling
	memstats, err := testutil.StopEnvoyMemoryProfile()
	require.NoError(t, err)

	require.GreaterOrEqualf(t, results.ActualQPS, *qps, "Actual QPS should be higher than target QPS")

	// TODO(musaprg): Currently, we're observing memory usage of envoy for checking memory stability,
	// but we need to use customized tinygo which enable to hook the timing of GC execution for finding
	// a best latency-friendly timing to invoke GC. Currently, tinygo invokes GC when memory allocations
	// attempts. It might affect the response latency. Basically, we don't need to invoke GC manually,
	// but it's better to do it out of the request-processing time.

	// Summarizing memory profile
	heapSizes := []float64{}
	allocSizes := []float64{}
	maxUsage := float64(0)
	maxAllocSize := float64(0)
	maxIndex := 0
	for i, m := range memstats {
		heapUsage := float64(m.Allocated) / float64(m.HeapSize)
		allocSize := float64(m.Allocated)
		heapSize := float64(m.HeapSize)
		allocSizes = append(allocSizes, allocSize)
		heapSizes = append(heapSizes, heapSize)
		if maxUsage < heapUsage {
			maxUsage = heapUsage
			maxIndex = i
		}
		if maxAllocSize < allocSize {
			maxAllocSize = allocSize
		}
	}
	log.Printf("peak memory usage: %v (elapsed %f sec after invoking load test)", maxUsage, float64(maxIndex*100)/1000)
	log.Printf("peak memory: %d bytes (+%d bytes increased from beginning)", int64(maxAllocSize), int64(maxAllocSize-allocSizes[0]))

	// Save the plot
	if *memoryUsageGraphDst != "" {
		envoyLog := stdErr.String()
		memStats, err := parseRuntimeMemStat(envoyLog)
		require.NoError(t, err, "Failed to parse memory stats", envoyLog)
		err = saveMemoryUsageGraph(memStats, *memoryUsageGraphDst)
		require.NoErrorf(t, err, "failed to save memory usage graph to %s", *memoryUsageGraphDst)
	}

	fortioLog.WriteTo(log.Writer())

	successRate := float64(results.RetCodes[200]) / float64(results.DurationHistogram.Count)
	require.GreaterOrEqual(t, successRate, targetSuccessRate, stdErr.String())
	require.LessOrEqual(t, results.DurationHistogram.Percentiles[0].Value, float64(targetNintyninthPercentileLatencyLimit), stdErr.String())
	require.NoErrorf(t, err, stdErr.String())
}

func saveMemoryUsageGraph(memStats runtimeMemStats, dst string) error {
	if dst == "" {
		return nil
	}

	// Plotting memory profile
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Heap profiling of envoy process (%f QPS, %s)", *qps, *targetExample)
	p.X.Label.Text = "elapsed time [ms]"
	p.Y.Label.Text = "memory size [KB]"
	heapSizePlot := make(plotter.XYs, len(memStats))
	allocSizePlot := make(plotter.XYs, len(memStats))
	for i, v := range memStats {
		t := (v.UnixNanoTime - memStats[0].UnixNanoTime) / 1000 // Convert to ms
		heapSizePlot[i].X = float64(t)
		heapSizePlot[i].Y = float64(v.HeapSize) / 1024 // Convert to KB
		allocSizePlot[i].X = float64(t)
		allocSizePlot[i].Y = float64(v.ReservedSize) / 1024 // Convert to KB
	}
	if err := plotutil.AddLinePoints(p,
		"heap_size", heapSizePlot,
		"allocated", allocSizePlot); err != nil {
		return err
	}

	if err := p.Save(font.Length(len(memStats))*vg.Millimeter, 8*vg.Inch, dst); err != nil {
		return err
	}

	return nil
}

type runtimeMemStat struct {
	UnixNanoTime int64  // ns
	HeapSize     uint64 // bytes
	ReservedSize uint64 // bytes
}

type runtimeMemStats []runtimeMemStat

func (s runtimeMemStats) Len() int {
	return len(s)
}
func (s runtimeMemStats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s runtimeMemStats) Less(i, j int) bool {
	return s[i].UnixNanoTime < s[j].UnixNanoTime
}

/// parseRuntimeMemStat parses the log of envoy and returns the tinygo's runtime GC stats.
func parseRuntimeMemStat(logs string) (runtimeMemStats, error) {
	var gcStats runtimeMemStats
	for _, result := range gcStatLogFormat.FindAllSubmatch([]byte(logs), -1) {
		if len(result) != 5 {
			return nil, fmt.Errorf("invalid memstat log format")
		}
		unixNanoTime, _ := strconv.ParseInt(string(result[2]), 10, 64)
		heapSize, _ := strconv.ParseUint(string(result[3]), 10, 64)
		reservedSize, _ := strconv.ParseUint(string(result[4]), 10, 64)
		gcStats = append(gcStats, runtimeMemStat{unixNanoTime, heapSize, reservedSize})
	}
	sort.Sort(gcStats)
	return gcStats, nil
}
