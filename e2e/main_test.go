// Copyright 2024 Tetrate
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
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

func checkMessage(str string, exps []string) bool {
	for _, exp := range exps {
		if !strings.Contains(str, exp) {
			return false
		}
	}
	return true
}

func startEnvoy(t *testing.T, adminPort int) (stdErr *bytes.Buffer, kill func()) {
	name := strings.TrimPrefix(t.Name(), "Test_")
	cmd := exec.Command("envoy",
		"--base-id", strconv.Itoa(adminPort),
		"--concurrency", "1", "--component-log-level", "wasm:trace",
		"-c", fmt.Sprintf("./examples/%s/envoy.yaml", name))

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())
	if !assert.Eventually(t, func() bool {
		res, err := http.Get(fmt.Sprintf("http://localhost:%d/listeners", adminPort))
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Envoy has not started") {
		t.Fatalf("Envoy stderr: %s", buf.String())
	}
	return buf, func() { require.NoError(t, cmd.Process.Kill()) }
}
