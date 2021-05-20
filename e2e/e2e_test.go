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

func (e *envoyPorts) getAdminAddress() string {
	return fmt.Sprintf("http://localhost:%d", e.admin)
}

func (e *envoyPorts) getEndpointAddress() string {
	return fmt.Sprintf("http://localhost:%d", e.endpoint)
}

func checkMessage(str string, exps, nexps []string) bool {
	for _, exp := range exps {
		if !strings.Contains(str, exp) {
			return false
		}
	}
	for _, nexp := range nexps {
		if strings.Contains(str, nexp) {
			return false
		}
	}
	return true
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

func testRunnerGetter(ps envoyPorts, r func(t *testing.T, nps envoyPorts, stdErr *bytes.Buffer)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		cmd, buf, conf := startEnvoy(t, ps)
		defer func() {
			require.NoError(t, cmd.Process.Kill())
			require.NoError(t, os.Remove(conf))
		}()
		r(t, ps, buf)
	}
}

func startEnvoy(t *testing.T, ps envoyPorts) (cmd *exec.Cmd, stdErr *bytes.Buffer, configPath string) {
	name := strings.TrimPrefix(t.Name(), "Test_E2E/")
	conf, err := getEnvoyConfigurationPath(t, name, ps)
	require.NoError(t, err)
	cmd = exec.Command("envoy",
		"--base-id", strconv.Itoa(ps.admin),
		"--concurrency", "1", "--component-log-level", "wasm:trace",
		"-c", conf)

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getAdminAddress() + "/listeners?format=json")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Envoy has not started")
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
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"helloworld: proxy_on_vm_start from Go!",
			"helloworld: It's",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func httpRouting(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var primary, canary bool
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress())
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

func httpAuthRandom(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	key := "this-is-key"
	value := "this-is-value"
	req, err := http.NewRequest("GET", ps.getEndpointAddress()+"/uuid", nil)
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

func httpHeaders(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", ps.getEndpointAddress(), nil)
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

func httpBody(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	req, err := http.NewRequest("GET", ps.getEndpointAddress()+"/anything",
		bytes.NewBuffer([]byte(`{ "example": "body" }`)))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		return checkMessage(stdErr.String(), []string{
			"body size: 21",
			`initial request body: { "example": "body" }`,
			"on http request body finished"},
			[]string{"failed to set request body", "failed to get request body"},
		) && checkMessage(string(body), []string{`"another": "body"`}, nil)
	}, 5*time.Second, 100*time.Millisecond, stdErr.String())
}

func network(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	key := "This-Is-Key"
	value := "this-is-value"
	req, err := http.NewRequest("GET", ps.getEndpointAddress(), nil)
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

func metrics(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var count int
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress())
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
		res, err := http.Get(ps.getAdminAddress() + "/stats")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		raw, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		return checkMessage(string(raw), []string{fmt.Sprintf("proxy_wasm_go.request_counter: %d", count)}, nil)
	}, 5*time.Second, time.Millisecond, "Expected stats not found")
}

func sharedData(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var count int
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress())
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

	out := stdErr.String()
	require.Contains(t, out, fmt.Sprintf("shared value: %d", count), out)
}

func sharedQueue(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress())
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy.")
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"dequeued data: hello",
			"dequeued data: world",
			"dequeued data: proxy-wasm",
		}, nil)
	}, 5*time.Second, time.Millisecond)
}

func vmPluginConfiguration(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	require.Eventually(t, func() bool {
		return checkMessage(stdErr.String(), []string{
			"name\": \"vm configuration", "name\": \"plugin configuration",
		}, nil)
	}, 5*time.Second, time.Millisecond, stdErr.String())
}

func configurationFromRoot(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress())
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

func accessLogger(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	exp := "/this/is/my/path"
	require.Eventually(t, func() bool {
		res, err := http.Get(ps.getEndpointAddress() + exp)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, time.Millisecond, "Endpoint not healthy")
	out := stdErr.String()
	require.Contains(t, out, exp, out)
}

func dispatchCallOnTick(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var count int = 1
	require.Eventually(t, func() bool {
		if strings.Contains(stdErr.String(), fmt.Sprintf("called! %d", count)) {
			count++
		}
		return count == 6
	}, 5*time.Second, 10*time.Millisecond, stdErr.String())
}

func callForeignOnTick(t *testing.T, ps envoyPorts, stdErr *bytes.Buffer) {
	var count int = 1
	require.Eventually(t, func() bool {
		if strings.Contains(stdErr.String(), fmt.Sprintf("foreign function (compress) called: %d", count)) {
			count++
		}
		return count == 6
	}, 5*time.Second, 10*time.Millisecond, stdErr.String())
}
