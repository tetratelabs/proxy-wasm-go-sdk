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
	"log"
	"runtime"
	"testing"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/fnet"
	"github.com/stretchr/testify/require"

	"github.com/guptarohit/asciigraph"
	"github.com/tetratelabs/proxy-wasm-go-sdk/e2e/testutil"
)

const (
	targetNintyninthPercentileLatencyLimit = 200 // ms
	targetSuccessRate                      = 0.7
)

var (
	qps           = flag.Float64("qps", 0, "QPS to run load test")
	duration      = flag.Int("duration", 10, "Duration of test in seconds")
	payloadSize   = flag.Int("payloadSize", 256, "Payload size in kilo bytes")
	targetExample = flag.String("targetExample", "http_headers", "Target example to run load test")
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
	log.Printf("\tDuration = %d [s], payloadSize = %d [byte]\n", *duration, *payloadSize)
	opts.QPS = *qps * 1.5 // In order to reach the target QPS, we need to set a little bit higer target QPS.
	opts.Duration = time.Duration(*duration) * time.Second

	// Run memory profiling to find out memory stability of SDK
	err := testutil.StartEnvoyMemoryProfile(8001)
	require.NoError(t, err)

	results, err := fhttp.RunHTTPTest(&opts)
	require.NoError(t, err)
	memstats, err := testutil.StopEnvoyMemoryProfile()
	require.NoError(t, err)
	require.GreaterOrEqualf(t, results.ActualQPS, *qps, "Actual QPS should be higher than target QPS")

	heapUsages := []float64{}
	allocSizes := []float64{}
	maxUsage := float64(0)
	maxIndex := 0
	for i, m := range memstats {
		heapUsage := float64(m.Allocated) / float64(m.HeapSize)
		heapUsages = append(heapUsages, heapUsage)
		allocSizes = append(allocSizes, float64(m.Allocated))
		if maxUsage < heapUsage {
			maxUsage = heapUsage
			maxIndex = i
		}
	}

	log.Printf("max memory usage: %v (elapsed %f sec after invoking load test)", maxUsage, float64(maxIndex*100)/1000)
	graph := asciigraph.Plot(allocSizes, asciigraph.Height(50))
	log.Printf("\n%s\n",graph)

	fortioLog.WriteTo(log.Writer())

	successRate := float64(results.RetCodes[200]) / float64(results.DurationHistogram.Count)
	require.GreaterOrEqual(t, successRate, targetSuccessRate, stdErr.String())
	require.LessOrEqual(t, results.DurationHistogram.Percentiles[0].Value, float64(targetNintyninthPercentileLatencyLimit), stdErr.String())
	require.NoErrorf(t, err, stdErr.String())
}
