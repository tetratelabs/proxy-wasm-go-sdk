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
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

const (
	exampleDefaultEndpointPort    = "18000"
	exampleDefaultStaticReplyPort = "8099"
	exampleDefaultAdminEndpoint   = "8001"
)

type envoyPorts struct {
	endpoint, staticReply, admin int
}

func Test_E2E(t *testing.T) {
	t.Run("network", testRunnerGetter(envoyPorts{
		endpoint:    11000,
		staticReply: 8000,
		admin:       28300,
	}, network))
	t.Run("shared_queue", testRunnerGetter(envoyPorts{
		endpoint:    11001,
		staticReply: 8001,
		admin:       28301,
	}, sharedQueue))
	t.Run("http_auth_random", testRunnerGetter(envoyPorts{
		endpoint:    11002,
		staticReply: 8002,
		admin:       28302,
	}, httpAuthRandom))
	t.Run("http_headers", testRunnerGetter(envoyPorts{
		endpoint:    11003,
		staticReply: 8003,
		admin:       28303,
	}, httpHeaders))
	t.Run("metrics", testRunnerGetter(envoyPorts{
		endpoint:    11004,
		staticReply: 8004,
		admin:       28304,
	}, metrics))
	t.Run("helloworld", testRunnerGetter(envoyPorts{
		endpoint:    11005,
		staticReply: 8005,
		admin:       28305,
	}, helloworld))
	t.Run("vm_plugin_configuration", testRunnerGetter(envoyPorts{
		endpoint:    11006,
		staticReply: 8006,
		admin:       28306,
	}, vmPluginConfiguration))
	t.Run("shared_data", testRunnerGetter(envoyPorts{
		endpoint:    11007,
		staticReply: 8007,
		admin:       28307,
	}, sharedData))
	t.Run("http_body", testRunnerGetter(envoyPorts{
		endpoint:    11008,
		staticReply: 8008,
		admin:       28308,
	}, httpBody))
	t.Run("configuration_from_root", testRunnerGetter(envoyPorts{
		endpoint:    11009,
		staticReply: 8009,
		admin:       28309,
	}, configurationFromRoot))
	t.Run("access_logger", testRunnerGetter(envoyPorts{
		endpoint:    11010,
		staticReply: 8010,
		admin:       28310,
	}, accessLogger))
	t.Run("dispatch_call_on_tick", testRunnerGetter(envoyPorts{
		endpoint:    11011,
		staticReply: 8011,
		admin:       28311,
	}, dispatchCallOnTick))
	t.Run("http_routing", testRunnerGetter(envoyPorts{
		endpoint:    11012,
		staticReply: 8012,
		admin:       28312,
	}, httpRouting))
	t.Run("foreign_call_on_tick", testRunnerGetter(envoyPorts{
		endpoint:    11013,
		staticReply: 8013,
		admin:       28313,
	}, callForeignOnTick))
}

type runner = func(t *testing.T, nps envoyPorts, stdErr *bytes.Buffer)

func testRunnerGetter(ps envoyPorts, r runner) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		cmd, buf, conf := startEnvoy(t, ps)
		r(t, ps, buf)
		defer func() {
			require.NoError(t, cmd.Process.Kill())
			require.NoError(t, os.Remove(conf))
		}()
	}
}

func startEnvoy(t *testing.T, ps envoyPorts) (cmd *exec.Cmd, stdErr *bytes.Buffer, configPath string) {
	name := strings.TrimPrefix(t.Name(), "Test_E2E/")
	conf, err := getEnvoyConfigurationPath(t, name, ps)
	require.NoError(t, err)
	cmd = exec.Command("envoy",
		"--base-id", strconv.Itoa(ps.admin),
		"--concurrency", "1",
		"-c", conf)

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second * 5)
	return cmd, buf, conf
}

func getEnvoyConfigurationPath(t *testing.T, name string, ps envoyPorts) (string, error) {
	bs, err := ioutil.ReadFile(fmt.Sprintf("./examples/%s/envoy.yaml", name))
	require.NoError(t, err)

	ms := strings.ReplaceAll(string(bs), exampleDefaultEndpointPort, strconv.Itoa(ps.endpoint))
	ms = strings.ReplaceAll(ms, exampleDefaultAdminEndpoint, strconv.Itoa(ps.admin))
	ms = strings.ReplaceAll(ms, exampleDefaultStaticReplyPort, strconv.Itoa(ps.staticReply))
	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.yaml")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(ms)
	require.NoError(t, err)
	return tmpFile.Name(), nil
}

func helloworld(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "helloworld: proxy_on_vm_start from Go!")
	require.Contains(t, out, "helloworld: It's")
}

func httpRouting(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var primary, canary bool
	for i := 0; i < 25; i++ { // TODO: maybe flaky
		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
		require.NoError(t, err)

		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		raw, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		body := string(raw)
		if strings.Contains(body, "canary") {
			canary = true
		}
		if strings.Contains(body, "primary") {
			primary = true
		}
		r.Body.Close()
		fmt.Println("received body: ", body)
	}

	out := stdErr.String()
	fmt.Println(out)
	require.True(t, primary, "must be routed to primary at least once")
	require.True(t, canary, "must be routed to canary at least once")
}

func httpAuthRandom(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	key := "this-is-key"
	value := "this-is-value"

	for i := 0; i < 25; i++ { // TODO: maybe flaky
		req, err := http.NewRequest("GET",
			fmt.Sprintf("http://localhost:%d/uuid", ps.endpoint), nil)
		require.NoError(t, err)
		req.Header.Add(key, value)

		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "access forbidden")
	require.Contains(t, out, "access granted")
	require.Contains(t, out, "response header from httpbin: :status: 200")
}

func httpHeaders(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	key := "this-is-key"
	value := "this-is-value"
	req.Header.Add(key, value)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, key)
	require.Contains(t, out, value)
	require.Contains(t, out, "server: envoy")
}

func httpBody(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/anything", ps.endpoint),
		bytes.NewBuffer([]byte(`{ "example": "body" }`)))
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "body size: 21")
	require.Contains(t, out, `initial request body: { "example": "body" }`)
	require.Contains(t, out, "on http request body finished")
	require.NotContains(t, out, "failed to set request body")
	require.NotContains(t, out, "failed to get request body")

	body, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), `"another": "body"`)
}

func network(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	key := "This-Is-Key"
	value := "this-is-value"

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	req.Header.Add(key, value)
	req.Header.Add("Connection", "close")

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	r.Body.Close()

	time.Sleep(time.Second * 5)

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, key)
	require.Contains(t, out, value)
	require.Contains(t, out, "downstream data received")
	require.Contains(t, out, "new connection!")
	require.Contains(t, out, "downstream connection close!")
	require.Contains(t, out, "upstream data received")
	require.Contains(t, out, "connection complete!")
	require.Contains(t, out, "remote address: 127.0.0.1:")
}

func metrics(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	fmt.Println(stdErr.String())

	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/stats", ps.admin), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	require.Contains(t, string(b), fmt.Sprintf("proxy_wasm_go.request_counter: %d", count))
}

func sharedData(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, fmt.Sprintf("shared value: %d", count))
}

func sharedQueue(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	time.Sleep(time.Second * 5)

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "dequeued data: hello")
	require.Contains(t, out, "dequeued data: world")
	require.Contains(t, out, "dequeued data: proxy-wasm")
}

func vmPluginConfiguration(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "name\": \"vm configuration")
	require.Contains(t, out, "name\": \"plugin configuration")
}

func configurationFromRoot(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ps.endpoint), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, "plugin config from root context")
	require.Contains(t, out, "name\": \"plugin configuration")
}

func accessLogger(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	exp := "/this/is/my/path"
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", ps.endpoint, exp), nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	require.Contains(t, out, exp)
}

func dispatchCallOnTick(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	time.Sleep(5 * time.Second)
	out := stdErr.String()
	fmt.Println(out)
	for i := 1; i < 6; i++ {
		require.Contains(t, out, fmt.Sprintf("called! %d", i))
	}
}

func callForeignOnTick(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	time.Sleep(5 * time.Second)
	out := stdErr.String()
	fmt.Println(out)
	for i := 1; i < 6; i++ {
		require.Contains(t, out, fmt.Sprintf("CallForeignFunction callNum: %d", i))
	}
}
