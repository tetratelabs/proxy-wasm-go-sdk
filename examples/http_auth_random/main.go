package main

import (
	"hash/fnv"
	"strconv"

	"github.com/mathetake/proxy-wasm-go/proxywasm"
	"github.com/mathetake/proxy-wasm-go/proxywasm/types"
)

func main() {
	proxywasm.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override default
func (ctx *httpHeaders) OnHttpRequestHeaders(int, bool) types.Action {
	hs, err := proxywasm.HostCallGetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCritical("failed to get request headers: ", err.Error())
		return types.ActionContinue
	}
	for _, h := range hs {
		proxywasm.LogInfo("request header: ", h[0]+": ", h[1])
	}

	if _, err := proxywasm.HostCallDispatchHttpCall(
		"httpbin", hs, "", [][2]string{}, 50000); err != nil {
		proxywasm.LogCritical("dipatch httpcall failed: ", err.Error())
	}

	return types.ActionPause
}

// override default
func (ctx *httpHeaders) OnHttpCallResponse(_ uint32, _ int, bodySize int, _ int) {
	b, err := proxywasm.HostCallGetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogCritical("failed to get response body: ", err.Error())
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	s := fnv.New32a()
	if _, err := s.Write(b); err != nil {
		proxywasm.LogCritical("failed to calculate hash: ", err.Error())
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	if s.Sum32()%2 == 0 {
		proxywasm.LogInfo("access granted")
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	msg := "access forbidden"
	proxywasm.LogInfo(msg)
	proxywasm.HostCallSendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go!!"},
	}, msg)
}

// override default
func (ctx *httpHeaders) OnLog() {
	proxywasm.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10), " finished")
}
