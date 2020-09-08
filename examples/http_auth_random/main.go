package main

import (
	"hash/fnv"
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
	"github.com/mathetake/proxy-wasm-go/runtime/types"
)

func main() {
	runtime.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	runtime.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) runtime.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override default
func (ctx *httpHeaders) OnHttpRequestHeaders(_ int, _ bool) types.Action {
	hs, err := runtime.HostCallGetHttpRequestHeaders()
	if err != nil {
		runtime.LogCritical("failed to get request headers: " + err.Error())
		return types.ActionContinue
	}
	for _, h := range hs {
		runtime.LogInfo("request header: " + h[0] + ": " + h[1])
	}

	if _, err := runtime.HostCallDispatchHttpCall(
		"httpbin", hs, "", [][2]string{}, 50000); err != nil {
		runtime.LogCritical("dipatch httpcall failed: " + err.Error())
	}

	return types.ActionPause
}

// override default
func (ctx *httpHeaders) OnHttpCallResponse(_ uint32, _ int, bodySize int, _ int) {
	b, err := runtime.HostCallGetHttpCallResponseBody(0, bodySize)
	if err != nil {
		runtime.LogCritical("failed to get response body: " + err.Error())
		runtime.HostCallResumeHttpRequest()
		return
	}

	s := fnv.New32a()
	if _, err := s.Write(b); err != nil {
		runtime.LogCritical("failed to calculate hash: " + err.Error())
		runtime.HostCallResumeHttpRequest()
		return
	}

	if s.Sum32()%2 == 0 {
		runtime.LogInfo("access granted")
		runtime.HostCallResumeHttpRequest()
		return
	}

	msg := "access forbidden"
	runtime.LogInfo(msg)
	runtime.HostCallSendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go!!"},
	}, msg)
}

// override default
func (ctx *httpHeaders) OnLog() {
	runtime.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10) + " finished")
}
