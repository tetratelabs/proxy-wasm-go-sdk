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

	"github.com/tetratelabs/proxy-wasm-go-sdk/e2e"
)

const (
	targetNintyninthPercentileLatencyLimit = 200 // ms
	targetSuccessRate                      = 0.8
)

var (
	qps           = flag.Float64("qps", 0, "QPS to run load test")
	duration      = flag.Int("duration", 10, "Duration of test in seconds")
	payloadSize   = flag.Int("payloadSize", 256, "Payload size in kilo bytes")
	targetExample = flag.String("targetExample", "http_headers", "Target example to run load test")
)

func Test_http_load(t *testing.T) {
	stdErr, kill, pss := e2e.StartEnvoyWith(*targetExample, t, 8001)
	defer kill()

	//initialMemoryStat := e2e.EnvoyMemoryUsage(t, 8001)
	initialMemorySize := pss()

	opts := fhttp.HTTPRunnerOptions{}
	opts.URL = "http://localhost:18000/uuid"
	opts.AllowInitialErrors = true
	opts.NumThreads = runtime.NumCPU()
	opts.NumConnections = (int(*qps) * *duration) / 2
	opts.Percentiles = []float64{99.0}

	fnet.ChangeMaxPayloadSize(*payloadSize * fnet.KILOBYTE)
	opts.Payload = fnet.Payload

	fortioLog := new(bytes.Buffer)
	opts.Out = fortioLog

	opts.HTTPReqTimeOut = 5000 * time.Second
	log.Printf("\tDuration = %d [s], payloadSize = %d [byte]\n", *duration, *payloadSize)
	opts.QPS = *qps
	opts.Duration = time.Duration(*duration) * time.Second
	results, err := fhttp.RunHTTPTest(&opts)

	//finalMemoryStat := e2e.EnvoyMemoryUsage(t, 8001)
	finalMemorySize := pss()

	log.Printf("\t\ttarget QPS: %v\n", opts.QPS)
	log.Printf("\t\tactual QPS: %v\n", results.ActualQPS)
	log.Printf("\tinitial memory size (PSS): %d bytes\n", initialMemorySize)
	log.Printf("\tfinal memory size (PSS): %d bytes (increased %f%%)\n", finalMemorySize, float64(finalMemorySize-initialMemorySize)/float64(initialMemorySize)*100)
	fortioLog.WriteTo(log.Writer())
	successRate := float64(results.RetCodes[200]) / float64(results.DurationHistogram.Count)
	require.GreaterOrEqual(t, successRate, targetSuccessRate, stdErr.String())
	require.LessOrEqual(t, results.DurationHistogram.Percentiles[0].Value, float64(targetNintyninthPercentileLatencyLimit), stdErr.String())
	require.NoErrorf(t, err, stdErr.String())
}
