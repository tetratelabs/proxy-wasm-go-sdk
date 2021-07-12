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
	"testing"

	"fortio.org/fortio/fhttp"
	"fortio.org/fortio/fnet"
	"github.com/stretchr/testify/require"
)

func Test_http_load(t *testing.T) {
	stdErr, kill := startEnvoyWith("network", t, 8001)
	defer kill()

	fnet.ChangeMaxPayloadSize(fnet.KILOBYTE)
	opts := fhttp.HTTPRunnerOptions{}
	numCalls := 100
	opts.Exactly = int64(numCalls)
	opts.QPS = float64(numCalls)
	opts.URL = "http://localhost:18000"
	fortioLog := new(bytes.Buffer)
	opts.Out = fortioLog
	_, err := fhttp.RunHTTPTest(&opts) // warm up round
	require.NoErrorf(t, err, stdErr.String(), fortioLog.String())

	megaByte := 1024 * fnet.KILOBYTE
	fnet.ChangeMaxPayloadSize(32 * megaByte)
	opts.Payload = fnet.Payload
	_, err = fhttp.RunHTTPTest(&opts)
	require.NoErrorf(t, err, stdErr.String(), fortioLog.String())
}
