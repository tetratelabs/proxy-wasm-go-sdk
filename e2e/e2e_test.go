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
	"io/ioutil"
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
			":status: 200", ":status: 503",
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

	for _, mode := range []string{
		"request",
		"response",
	} {
		t.Run(mode, func(t *testing.T) {
			for _, tc := range []struct {
				op, expBody string
			}{
				{op: "append", expBody: `[original body][this is appended body]`},
				{op: "prepend", expBody: `[this is prepended body][original body]`},
				{op: "replace", expBody: `[this is replaced body]`},
				// Should fall back to to the replace.
				{op: "invalid", expBody: `[this is replaced body]`},
			} {
				tc := tc
				t.Run(tc.op, func(t *testing.T) {
					require.Eventually(t, func() bool {
						req, err := http.NewRequest("PUT", "http://localhost:18000/anything",
							bytes.NewBuffer([]byte(`[original body]`)))
						require.NoError(t, err)
						req.Header.Add("buffer-replace-at", mode)
						req.Header.Add("buffer-operation", tc.op)
						res, err := http.DefaultClient.Do(req)
						if err != nil {
							return false
						}
						defer res.Body.Close()
						body, err := io.ReadAll(res.Body)
						require.NoError(t, err)
						require.Equal(t, tc.expBody, string(body))
						require.True(t, checkMessage(stdErr.String(), []string{
							fmt.Sprintf(`original %s body: [original body]`, mode)},
							[]string{"failed to"},
						))
						return true
					}, 10*time.Second, 500*time.Millisecond, stdErr.String())
				})
			}
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

	const customHeaderKey = "my-custom-header"
	customHeaderToExpectedCounts := map[string]int{
		"foo": 3,
		"bar": 5,
	}
	for headerValue, expCount := range customHeaderToExpectedCounts {
		var actualCount int
		require.Eventually(t, func() bool {
			req, err := http.NewRequest("GET", "http://localhost:18000", nil)
			require.NoError(t, err)
			req.Header.Add(customHeaderKey, headerValue)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return false
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				return false
			}
			actualCount++
			return actualCount == expCount
		}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	}

	for headerValue, expCount := range customHeaderToExpectedCounts {
		expectedMetric := fmt.Sprintf("custom_header_value_counts{value=\"%s\",reporter=\"wasmgosdk\"} %d", headerValue, expCount)
		require.Eventually(t, func() bool {
			res, err := http.Get("http://localhost:8001/stats/prometheus")
			if err != nil {
				return false
			}
			defer res.Body.Close()
			raw, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			return checkMessage(string(raw), []string{expectedMetric}, nil)
		}, 5*time.Second, time.Millisecond, "Expected stats not found")
	}
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
			"upsteam cluster matadata location[region]=ap-northeast-1",
			"upsteam cluster matadata location[cloud_provider]=aws",
			"upsteam cluster matadata location[az]=ap-northeast-1a",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_postpone_requests(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return checkMessage(stdErr.String(), []string{
			"postpone request with contextID=2",
			"resume request with contextID=2",
		}, nil)
	}, 6*time.Second, time.Millisecond, stdErr.String())
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

func Test_json_validation(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()

	require.Eventually(t, func() bool {
		req, _ := http.NewRequest("GET", "http://localhost:18000", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()

		_, err = io.Copy(ioutil.Discard, res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, res.StatusCode)

		jsonBody := `{"id": "abc123", "token": "xyz456"}`

		req, _ = http.NewRequest("POST", "http://localhost:18000", strings.NewReader(jsonBody))
		req.Header.Add("Content-Type", "application/json")
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()

		_, err = io.Copy(ioutil.Discard, res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		return true
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_multiple_dispatches(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()

	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil || res.StatusCode != http.StatusOK {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")

	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"wasm log: pending dispatched requests: 9",
			"wasm log: pending dispatched requests: 8",
			"wasm log: pending dispatched requests: 7",
			"wasm log: pending dispatched requests: 6",
			"wasm log: pending dispatched requests: 5",
			"wasm log: pending dispatched requests: 4",
			"wasm log: pending dispatched requests: 3",
			"wasm log: pending dispatched requests: 2",
			"wasm log: pending dispatched requests: 1",
			"wasm log: response resumed after processed 10 dispatched request",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}
