// Copyright 2020 Tetrate
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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

const (
	envoyVersion       = "wasm:1.15"
	envoyEndpoint      = "http://localhost:18000"
	envoyAdminEndpoint = "http://localhost:8001"
)

func startExample(t *testing.T, name string) (*exec.Cmd, *bytes.Buffer) {
	cmd := exec.Command("getenvoy",
		"run",
		envoyVersion,
		"--",
		"--concurrency", "2",
		"-c", fmt.Sprintf("./examples/%s/envoy.yaml", name))

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second * 5) // TODO: use admin endpoint to check health
	return cmd, buf
}

func TestE2E_helloworld(t *testing.T) {
	cmd, stdErr := startExample(t, "helloworld")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, "wasm log helloworld: proxy_on_vm_start from Go!"))
	assert.True(t, strings.Contains(out, "wasm log helloworld: OnTick on "))
}

func TestE2E_http_auth_random(t *testing.T) {
	cmd, stdErr := startExample(t, "http_auth_random")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	key := "this-is-key"
	value := "this-is-value"

	for i := 0; i < 25; i++ { // TODO: maybe flaky
		req, err := http.NewRequest("GET", envoyEndpoint+"/uuid", nil)
		require.NoError(t, err)
		req.Header.Add(key, value)

		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, "access forbidden"))
	assert.True(t, strings.Contains(out, "access granted"))
	assert.True(t, strings.Contains(out, "response header from httpbin: :status: 200"))
}

func TestE2E_http_headers(t *testing.T) {
	cmd, stdErr := startExample(t, "http_headers")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	req, err := http.NewRequest("GET", envoyEndpoint, nil)
	require.NoError(t, err)

	key := "this-is-key"
	value := "this-is-value"
	req.Header.Add(key, value)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, key))
	assert.True(t, strings.Contains(out, value))
	assert.True(t, strings.Contains(out, "server: envoy"))
}

func TestE2E_network(t *testing.T) {
	cmd, stdErr := startExample(t, "network")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	key := "This-Is-Key"
	value := "this-is-value"

	doReq := func() {
		req, err := http.NewRequest("GET", envoyEndpoint, nil)
		require.NoError(t, err)

		req.Header.Add(key, value)
		req.Header.Add("Connection", "close")

		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	doReq()

	time.Sleep(time.Second)

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, key))
	assert.True(t, strings.Contains(out, value))
	assert.True(t, strings.Contains(out, "downstream data received"))
	assert.True(t, strings.Contains(out, "new connection!"))
	assert.True(t, strings.Contains(out, "downstream connection close!"))
	assert.True(t, strings.Contains(out, "upstream data received"))
	assert.True(t, strings.Contains(out, "connection complete!"))
	assert.True(t, strings.Contains(out, "remote address: 127.0.0.1:8099"))
}

func TestE2E_metrics(t *testing.T) {
	cmd, stdErr := startExample(t, "metrics")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	req, err := http.NewRequest("GET", envoyEndpoint, nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	fmt.Println(stdErr.String())

	req, err = http.NewRequest("GET", envoyAdminEndpoint+"/stats", nil)
	require.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(b), fmt.Sprintf("proxy_wasm_go.request_counter: %d", count)))
}

func TestE2E_shared_data(t *testing.T) {
	cmd, stdErr := startExample(t, "shared_data")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	req, err := http.NewRequest("GET", envoyEndpoint, nil)
	require.NoError(t, err)

	count := 10
	for i := 0; i < count; i++ {
		r, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		r.Body.Close()
	}

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, fmt.Sprintf("shared value: %d", count)))
}

func TestE2E_shared_queue(t *testing.T) {
	cmd, stdErr := startExample(t, "shared_queue")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	req, err := http.NewRequest("GET", envoyEndpoint, nil)
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
	assert.True(t, strings.Contains(out, "dequeued data: hello"))
	assert.True(t, strings.Contains(out, "dequeued data: world"))
	assert.True(t, strings.Contains(out, "dequeued data: proxy-wasm"))
}

func TestE2E_vm_plugin_configuration(t *testing.T) {
	cmd, stdErr := startExample(t, "vm_plugin_configuration")
	defer func() {
		require.NoError(t, cmd.Process.Kill())
	}()

	out := stdErr.String()
	fmt.Println(out)
	assert.True(t, strings.Contains(out, "name\": \"vm configuration"))
	assert.True(t, strings.Contains(out, "name\": \"plugin configuration"))
}
