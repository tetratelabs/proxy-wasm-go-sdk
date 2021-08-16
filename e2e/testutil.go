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
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CheckMessage(str string, exps, nexps []string) bool {
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

// startEnvoyWith is used for invoking the envoy process with a specified example.
func StartEnvoyWith(name string, t *testing.T, adminPort int) (stdErr *bytes.Buffer, kill func()) {
	cmd := exec.Command("envoy",
		"--base-id", strconv.Itoa(adminPort),
		"--concurrency", "1", "--component-log-level", "wasm:trace",
		"-c", fmt.Sprintf("./examples/%s/envoy.yaml", name))

	buf := new(bytes.Buffer)
	cmd.Stderr = buf
	require.NoError(t, cmd.Start())
	require.Eventually(t, func() bool {
		res, err := http.Get(fmt.Sprintf("http://localhost:%d/listeners", adminPort))
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond, "Envoy has not started")
	return buf, func() { require.NoError(t, cmd.Process.Kill()) }
}

// startEnvoy is used for invoking the envoy process which is used for e2e testing.
// The target example is selected based on the name of the test case.
func StartEnvoy(t *testing.T, adminPort int) (stdErr *bytes.Buffer, kill func()) {
	name := strings.TrimPrefix(t.Name(), "Test_")
	return StartEnvoyWith(name, t, adminPort)
}

// MemoryStat represents the response format of :8081/memory which returns memory stat for envoy process.
type MemoryStat struct {
	Allocated memoryBytes `json:"allocated"`
	HeapSize memoryBytes `json:"heap_size"`
	PageheapUnmapped memoryBytes `json:"pageheap_unmapped"`
	PageheapFree memoryBytes `json:"pageheap_free"`
	TotalThreadCache memoryBytes `json:"total_thread_cache"`
	TotalPhysicalBytes memoryBytes `json:"total_physical_bytes"`
}

type memoryBytes int64

func (m *memoryBytes) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*m = memoryBytes(i)
	return nil
}

// EnvoyMemoryUsage is used for getting the memory usage of envoy process.
func EnvoyMemoryUsage(t *testing.T, adminPort int) MemoryStat {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/memory", adminPort))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, err)
	defer resp.Body.Close()

	var memory MemoryStat
	if err := json.NewDecoder(resp.Body).Decode(&memory); err != nil {
		require.NoError(t, err)
	}
	return memory
}
