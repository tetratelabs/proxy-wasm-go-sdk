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
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_access_logger(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	exp := "/this/is/my/path"
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000" + exp)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy")

	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{exp}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_configuration_from_root(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"plugin config from root context",
			"name\": \"plugin configuration",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_dispatch_call_on_tick(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	var count int = 1
	require.Eventually(t, func() bool {
		if strings.Contains(stdErr.String(), fmt.Sprintf("called! %d", count)) {
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
			"helloworld: proxy_on_vm_start from Go!",
			"helloworld: It's",
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
	req, err := http.NewRequest("GET", "http://localhost:18000/anything",
		bytes.NewBuffer([]byte(`{ "initial": "body" }`)))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		return checkMessage(stdErr.String(), []string{
			"body size: 21",
			`initial request body: { "initial": "body" }`,
			"on http request body finished"},
			[]string{"failed to set request body", "failed to get request body"},
		) && checkMessage(string(body), []string{`"another": "body"`}, nil)
	}, 5*time.Second, 500*time.Millisecond, stdErr.String())
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
		raw, err := ioutil.ReadAll(res.Body)
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
		res, err := http.Get("http://localhost:8001/stats")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		raw, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		return checkMessage(string(raw), []string{fmt.Sprintf("proxy_wasm_go.request_counter: %d", count)}, nil)
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
		return checkMessage(stdErr.String(), []string{fmt.Sprintf("shared value: %d", count)}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func Test_shared_queue(t *testing.T) {
	stdErr, kill := startEnvoy(t, 8001)
	defer kill()
	require.Eventually(t, func() bool {
		res, err := http.Get("http://localhost:18000")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			`enqueued data: {"key": ":method","value": "GET"}`,
			`dequeued data: {"key": ":method","value": "GET"}`,
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
