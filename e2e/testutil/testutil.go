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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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
// The fields are defined in https://www.envoyproxy.io/docs/envoy/latest/api-v3/admin/v3/memory.proto.html.
type MemoryStat struct {
	// The number of bytes allocated by the heap for envoy.
	Allocated memoryBytes `json:"allocated"`
	// The number of bytes reserved for the heap.
	HeapSize memoryBytes `json:"heap_size"`
	// The number of bytes in free, unmapped pages in the page heap.
	PageheapUnmapped memoryBytes `json:"pageheap_unmapped"`
	// The number of bytes in free, mapped pages in the page heap.
	PageheapFree memoryBytes `json:"pageheap_free"`
	// The amount of memory used by the TCMalloc thread caches.
	TotalThreadCache memoryBytes `json:"total_thread_cache"`
	// The number of bytes of the physical memory usage by the allocator.
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

var memoryStats struct {
	sync.Mutex
	memstats []MemoryStat
}

// EnvoyMemoryProfile represents a memory profile of envoy process.
type EnvoyMemoryProfile struct {
	mu        sync.Mutex
	profiling bool
	ctx       context.Context
	cancel    func()
	done      chan bool
}

var envoyMemoryProfile EnvoyMemoryProfile

// StartEnvoyMemoryProfile starts the memory profiling of envoy process.
func StartEnvoyMemoryProfile(adminPort int) error {
	envoyMemoryProfile.mu.Lock()
	defer envoyMemoryProfile.mu.Unlock()
	if envoyMemoryProfile.done == nil {
		envoyMemoryProfile.done = make(chan bool)
	}
	if envoyMemoryProfile.ctx == nil {
		ctx := context.Background()
		envoyMemoryProfile.ctx, envoyMemoryProfile.cancel = context.WithCancel(ctx)
	}
	if envoyMemoryProfile.profiling {
		return fmt.Errorf("profiling is already running")
	}
	envoyMemoryProfile.profiling = true
	go runProfile(envoyMemoryProfile.ctx, adminPort)
	return nil
}

// runProfile is profiling the memory usage of envoy process for every 100ms.
func runProfile(ctx context.Context, adminPort int) {
	var err error
	for {
		time.Sleep(100 * time.Millisecond)
		var m *MemoryStat
		m, e := EnvoyMemoryUsage(adminPort)
		if e != nil {
			err = e
		}
		memoryStats.Lock()
		memoryStats.memstats = append(memoryStats.memstats, *m)
		memoryStats.Unlock()
		select {
		case <-ctx.Done():
			if err != nil {
				log.Fatal(err)
			}
			envoyMemoryProfile.done <- true
		default:

		}
	}
}

// StopEnvoyMemoryProfile stops the memory profiling of envoy process.
func StopEnvoyMemoryProfile() ([]MemoryStat, error) {
	envoyMemoryProfile.mu.Lock()
	defer envoyMemoryProfile.mu.Unlock()
	if !envoyMemoryProfile.profiling {
		return nil, fmt.Errorf("profiling is not running")
	}
	envoyMemoryProfile.profiling = false
	envoyMemoryProfile.cancel()
	<-envoyMemoryProfile.done

	memoryStats.Lock()
	defer memoryStats.Unlock()
	memstats := memoryStats.memstats
	memoryStats.memstats = make([]MemoryStat, 0)
	return memstats, nil
}

// EnvoyMemoryUsage is used for getting the memory usage of envoy process.
func EnvoyMemoryUsage(adminPort int) (*MemoryStat, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/memory", adminPort))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var memory MemoryStat
	err = json.NewDecoder(resp.Body).Decode(&memory)
	return &memory, err
}
