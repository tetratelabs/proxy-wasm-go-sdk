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
	targetQPS = 100
)

func Test_http_load(t *testing.T) {
	stdErr, kill := e2e.StartEnvoyWith("network", t, 8001)
	defer kill()

	states := []struct {
		numCalls          int64
		payloadSize       int
		upperLimitLatency float64
	}{
		{100, 256 * fnet.KILOBYTE, 100},
		{1, 16384 * fnet.KILOBYTE, 150},
	}

	opts := fhttp.HTTPRunnerOptions{}
	opts.URL = "http://localhost:18000"
	opts.AllowInitialErrors = true
	opts.NumThreads = runtime.NumCPU()
	opts.Percentiles = []float64{99.0}

	fnet.ChangeMaxPayloadSize(fnet.KILOBYTE)
	opts.Payload = fnet.Payload

	fortioLog := new(bytes.Buffer)
	opts.Out = fortioLog

	opts.Exactly = 1
	_, err := fhttp.RunHTTPTest(&opts) // warm up round
	require.NoErrorf(t, err, stdErr.String(), fortioLog.String())

	opts.HTTPReqTimeOut = 5000 * time.Second
	opts.AbortOn = -1

	for _, state := range states {
		stdErr.Reset()
		fortioLog.Reset()
		log.Printf("\tnumCalls = %d, payloadSize = %d [byte]\n", state.numCalls, state.payloadSize)
		fnet.ChangeMaxPayloadSize(state.payloadSize)
		opts.Payload = fnet.Payload
		opts.Exactly = state.numCalls
		opts.QPS = float64(targetQPS)
		results, err := fhttp.RunHTTPTest(&opts)
		log.Printf("\t\ttarget QPS: %v\n", targetQPS)
		log.Printf("\t\tactual QPS: %v\n", results.ActualQPS)
		require.Equal(t, results.DurationHistogram.Count, results.RetCodes[200], stdErr.String(), fortioLog.String())
		require.LessOrEqual(t, results.DurationHistogram.Percentiles[0].Value, state.upperLimitLatency, stdErr.String(), fortioLog.String())
		require.NoErrorf(t, err, stdErr.String(), fortioLog.String())
	}
}
