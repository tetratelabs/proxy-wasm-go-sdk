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
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/proxy-wasm-go-sdk/e2e/testutil"
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func Test_E2E(t *testing.T) {
	t.Run("network", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11000,
		StaticReply: 8000,
		Admin:       28300,
	}, network))
	t.Run("shared_queue", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11001,
		StaticReply: 8001,
		Admin:       28301,
	}, sharedQueue))
	t.Run("http_auth_random", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11002,
		StaticReply: 8002,
		Admin:       28302,
	}, httpAuthRandom))
	t.Run("http_headers", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11003,
		StaticReply: 8003,
		Admin:       28303,
	}, httpHeaders))
	t.Run("metrics", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11004,
		StaticReply: 8004,
		Admin:       28304,
	}, metrics))
	t.Run("helloworld", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11005,
		StaticReply: 8005,
		Admin:       28305,
	}, helloworld))
	t.Run("vm_plugin_configuration", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11006,
		StaticReply: 8006,
		Admin:       28306,
	}, vmPluginConfiguration))
	t.Run("shared_data", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11007,
		StaticReply: 8007,
		Admin:       28307,
	}, sharedData))
	t.Run("http_body", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11008,
		StaticReply: 8008,
		Admin:       28308,
	}, httpBody))
	t.Run("configuration_from_root", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11009,
		StaticReply: 8009,
		Admin:       28309,
	}, configurationFromRoot))

	t.Run("access_logger", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11010,
		StaticReply: 8010,
		Admin:       28310,
	}, accessLogger))
	t.Run("dispatch_call_on_tick", testutil.TestRunnerGetter(testutil.EnvoyPorts{
		Endpoint:    11011,
		StaticReply: 8011,
		Admin:       28311,
	}, dispatchCallOnTick))
}

func helloworld(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "helloworld: proxy_on_vm_start from Go!")
	assert.Contains(t, out, "helloworld: It's")
}

func httpAuthRandom(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	key := "this-is-key"
	value := "this-is-value"

	for i := 0; i < 25; i++ { // TODO: maybe flaky
		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://localhost:%d/uuid", ps.Endpoint), nil)
		require.NoError(t, err)
		req.Header.Add(key, value)

		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "access forbidden")
	assert.Contains(t, out, "access granted")
	assert.Contains(t, out, "response header from httpbin: :status: 200")
}

func httpHeaders(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	key := "this-is-key"
	value := "this-is-value"
	req.Header.Add(key, value)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, key)
	assert.Contains(t, out, value)
	assert.Contains(t, out, "server: envoy")
}

func httpBody(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/anything", ps.Endpoint),
		bytes.NewBuffer([]byte(`{ "example": "body" }`)))
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "body size: 21")
	assert.Contains(t, out, `initial request body: { "example": "body" }`)
	assert.Contains(t, out, "on http request body finished")
	assert.NotContains(t, out, "failed to set request body")
	assert.NotContains(t, out, "failed to get request body")

	body, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"another": "body"`)
}

func network(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	key := "This-Is-Key"
	value := "this-is-value"

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	req.Header.Add(key, value)
	req.Header.Add("Connection", "close")

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	r.Body.Close()

	time.Sleep(time.Second)

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, key)
	assert.Contains(t, out, value)
	assert.Contains(t, out, "downstream data received")
	assert.Contains(t, out, "new connection!")
	assert.Contains(t, out, "downstream connection close!")
	assert.Contains(t, out, "upstream data received")
	assert.Contains(t, out, "connection complete!")
	assert.Contains(t, out, "remote address: 127.0.0.1:")
}

func metrics(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	fmt.Println(stdErr.String())

	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/stats", ps.Admin), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	assert.Contains(t, string(b), fmt.Sprintf("proxy_wasm_go.request_counter: %d", count))
}

func sharedData(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, fmt.Sprintf("shared value: %d", count))
}

func sharedQueue(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	time.Sleep(time.Second * 2)

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "dequeued data: hello")
	assert.Contains(t, out, "dequeued data: world")
	assert.Contains(t, out, "dequeued data: proxy-wasm")
}

func vmPluginConfiguration(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "name\": \"vm configuration")
	assert.Contains(t, out, "name\": \"plugin configuration")
}

func configurationFromRoot(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.Endpoint), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, "plugin config from root context")
	assert.Contains(t, out, "name\": \"plugin configuration")
}

func accessLogger(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	exp := "/this/is/my/path"
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", ps.Endpoint, exp), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	assert.Contains(t, out, exp)
}

func dispatchCallOnTick(t *testing.T, ps testutil.EnvoyPorts, stdErr *bytes.Buffer) {
	time.Sleep(3 * time.Second)
	out := stdErr.String()
	fmt.Println(out)
	for i := 1; i < 6; i++ {
		assert.Contains(t, out, fmt.Sprintf("called! %d", i))
	}
}
