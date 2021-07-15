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

package e2e

import (
	"bytes"
	"log"
	"testing"
	"time"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/fnet"
	"github.com/stretchr/testify/require"
)

func Test_http_load(t *testing.T) {
	stdErr, kill := startEnvoyWith("network", t, 8001)
	defer kill()

	states := []struct {
		numCalls    int64
		payloadSize int
	}{
		{1, 256 * fnet.KILOBYTE},
		{1, 512 * fnet.KILOBYTE},
		{1, 1024 * fnet.KILOBYTE},
		{1, 2048 * fnet.KILOBYTE},
		{1, 4096 * fnet.KILOBYTE},
		{1, 8192 * fnet.KILOBYTE},
		{1, 16384 * fnet.KILOBYTE},
		{1, 32798 * fnet.KILOBYTE},
	}

	opts := fhttp.HTTPRunnerOptions{}
	opts.URL = "http://localhost:18000"
	opts.AllowInitialErrors = true

	fnet.ChangeMaxPayloadSize(fnet.KILOBYTE)
	opts.Payload = fnet.Payload

	fortioLog := new(bytes.Buffer)
	opts.Out = fortioLog

	opts.Exactly = 100
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
		results, err := fhttp.RunHTTPTest(&opts)
		log.Printf("\tReturn Codes: %v\n", results.RetCodes)
		require.Equal(t, results.RetCodes[200], state.numCalls, stdErr.String(), fortioLog.String())
		require.NoErrorf(t, err, stdErr.String(), fortioLog.String())
	}
}
