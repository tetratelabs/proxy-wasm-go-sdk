// Copyright 2020-2021 Tetrate
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
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_dispatch_call_on_tick(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	var count int = 1
	require.Eventually(t, func() bool {
		if checkMessage(stdErr.String(), []string{
			fmt.Sprintf("called %d for contextID=1", count),
			fmt.Sprintf("called %d for contextID=2", count),
		}, nil) {
			count++
		}
		return count == 6
	}, 5*time.Second, 10*time.Millisecond, stdErr.String())
}

func Test_foreign_call_on_tick(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	var count int = 1
	require.Eventually(t, func() bool {
		if strings.Contains(stdErr.String(), fmt.Sprintf("foreign function (compress) called: %d", count)) {
			count++
		}
		return count == 6
	}, 5*time.Second, 10*time.Millisecond, stdErr.String())
}

func Test_helloworld(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"OnPluginStart from Go!",
			"It's",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_http_auth_random(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	key := "this-is-key"
	value := "this-is-value"
	req, err := http.NewRequest("GET", "http://localhost:18000/uuid", nil)
	require.NoError(t, err)
	req.Header.Add(key, value)
	require.Eventually(t, func() bool {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return checkMessage(stdErr.String(), []string{
			"access forbidden",
			"access granted",
			"response header from httpbin: :status: 200",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_http_body(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()

	for _, tc := range []struct {
		op, expBody string
	}{
		{op: "append", expBody: `[original body][this is appended body]`},
		{op: "prepend", expBody: `[this is prepended body][original body]`},
		{op: "replace", expBody: `[this is replaced body]`},
		// Shoud fall back to to the replace.
		{op: "invalid", expBody: `[this is replaced body]`},
	} {
		t.Run(tc.op, func(t *testing.T) {
			require.Eventually(t, func() bool {
				req, err := http.NewRequest("PUT", "http://localhost:18000/anything",
					bytes.NewBuffer([]byte(`[original body]`)))
				require.NoError(t, err)
				req.Header.Add("buffer-operation", tc.op)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return false
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				return string(body) == tc.expBody &&
					checkMessage(stdErr.String(), []string{
						`original request body: [original body]`},
						[]string{"failed to"},
					) && checkMessage(string(body), []string{tc.expBody}, nil)
			}, 5*time.Second, 500*time.Millisecond, stdErr.String())
		})
	}
}

func Test_http_headers(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	req, err := http.NewRequest("GET", "http://localhost:18000/uuid", nil)
	require.NoError(t, err)
	key := "this-is-key"
	value := "this-is-value"
	req.Header.Add(key, value)
	require.Eventually(t, func() bool {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return checkMessage(stdErr.String(), []string{
			key, value, "server: envoy",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_http_routing(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	var primary, canary bool
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		raw, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		defer res.Body.Close()
		body := string(raw)
		if strings.Contains(body, "canary") {
			canary = true
		}
		if strings.Contains(body, "primary") {
			primary = true
		}
		return primary && canary
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_metrics(t *testing.T) {
	_, kill := startEnvoy(t, 8001)
	defer kill()
	var count int
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return false
		}
		count++
		return count == 10
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:8001/stats/prometheus")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		raw, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		return checkMessage(string(raw), []string{fmt.Sprintf("proxy_wasm_go_request_counter{} %d", count)}, nil)
	}, 5*time.Second, time.Millisecond, "Expected stats not found")
}

func Test_network(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	key := "This-Is-Key"
	value := "this-is-value"
	req, err := http.NewRequest("GET", "http://localhost:18000", nil)
	require.NoError(t, err)
	req.Header.Add(key, value)
	req.Header.Add("Connection", "close")
	require.Eventually(t, func() bool {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return checkMessage(stdErr.String(), []string{
			key, value,
			"downstream data received",
			"new connection!",
			"downstream connection close!",
			"upstream data received",
			"connection complete!",
			"remote address: 127.0.0.1:",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_shared_data(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	var count int = 10000000
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return false
		}
		count++
		return count == 10000010
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{fmt.Sprintf("shared value: %d", count)}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_shared_queue(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil || res.StatusCode != http.StatusOK {
			return false
		}
		defer res.Body.Close()

		res, err = http.Get("http://localhost:18001")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			`enqueued data: {"key": ":method","value": "GET"}`,
			`dequeued data from http_request_headers`,
			`dequeued data from http_response_headers`,
			`dequeued data from tcp_data_hashes`,
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_vm_plugin_configuration(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"name\": \"vm configuration", "name\": \"plugin configuration",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}
