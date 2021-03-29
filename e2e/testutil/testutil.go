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

package testutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	ExampleDefaultEndpointPort    = "18000"
	ExampleDefaultStaticReplyPort = "8099"
	ExampleDefaultAdminEndpoint   = "8001"
)

type EnvoyPorts struct {
	Endpoint, StaticReply, Admin int
}

type runner = func(t *testing.T, nps EnvoyPorts, stdErr *bytes.Buffer)

func TestRunnerGetter(ps EnvoyPorts, r runner) func(t *testing.T) {
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

func startEnvoy(t *testing.T, ps EnvoyPorts) (cmd *exec.Cmd, stdErr *bytes.Buffer, configPath string) {
	name := strings.TrimPrefix(t.Name(), "Test_E2E/")
	conf, err := getEnvoyConfigurationPath(t, name, ps)
	require.NoError(t, err)
	cmd = exec.Command("envoy",
		"--base-id", strconv.Itoa(ps.Admin),
		"--concurrency", "1",
		"-c", conf)

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second * 5)
	return cmd, buf, conf
}

func getEnvoyConfigurationPath(t *testing.T, name string, ps EnvoyPorts) (string, error) {
	bs, err := ioutil.ReadFile(fmt.Sprintf("./examples/%s/envoy.yaml", name))
	require.NoError(t, err)

	ms := strings.ReplaceAll(string(bs), ExampleDefaultEndpointPort, strconv.Itoa(ps.Endpoint))
	ms = strings.ReplaceAll(ms, ExampleDefaultAdminEndpoint, strconv.Itoa(ps.Admin))
	ms = strings.ReplaceAll(ms, ExampleDefaultStaticReplyPort, strconv.Itoa(ps.StaticReply))
	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.yaml")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(ms)
	require.NoError(t, err)
	return tmpFile.Name(), nil
}
