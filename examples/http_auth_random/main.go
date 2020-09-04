package main

import (
	"hash/fnv"
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
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
func (ctx *httpHeaders) OnHttpRequestHeaders(_ int, _ bool) runtime.Action {
	hs, st := ctx.GetHttpRequestHeaders()
	if st != runtime.StatusOk {
		runtime.LogCritical("failed to get request headers")
		return runtime.ActionContinue
	}
	for _, h := range hs {
		runtime.LogInfo("request header: " + h[0] + ": " + h[1])
	}

	ctx.DispatchHttpCall("httpbin", hs, "", [][2]string{}, 50000)
	return runtime.ActionPause
}

// override default
func (ctx *httpHeaders) OnHttpCallResponse(_ uint32, _ int, bodySize int, _ int) {
	b, st := ctx.GetHttpCallResponseBody(0, bodySize)
	if st != runtime.StatusOk {
		runtime.LogCritical("failed to get response body")
		ctx.ResumeHttpRequest()
		return
	}

	s := fnv.New32a()
	if _, err := s.Write(b); err != nil {
		runtime.LogCritical("failed to calculate hash: " + err.Error())
		ctx.ResumeHttpRequest()
		return
	}

	if s.Sum32()%2 == 0 {
		runtime.LogInfo("access granted")
		ctx.ResumeHttpRequest()
		return
	}

	msg := "access forbidden"
	runtime.LogInfo(msg)
	ctx.SendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go!!"},
	}, msg)
}

// override default
func (ctx *httpHeaders) OnLog() {
	runtime.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10) + " finished")
}
