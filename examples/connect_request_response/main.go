package main

import (
	"encoding/binary"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"math/rand"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

const (
	sharedDataKey = "shared_data_key"
)

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("<---- 新Http连接 ---->")
	return &responseContext{contextID: contextID}
}

type responseContext struct {
	contextID uint32
	types.DefaultHttpContext
}

// Override types.DefaultHttpContext.
func (r *responseContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	/*headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		return 0
	}
	for i := 0; i < len(headers); i++ {
		for j := 0; j < len(headers[i]); j++ {
			proxywasm.LogInfof("请求头中headers arr[%v][%v] = %v", i, j, headers[i][j])
		}
	}*/

	initialValueBuf := make([]byte, 8)
	rand := rand.Uint64() * 100000000
	binary.LittleEndian.PutUint64(initialValueBuf, rand)
	proxywasm.SetSharedData(sharedDataKey, initialValueBuf, 0)
	proxywasm.LogInfof("请求头中的share data: %d", rand)

	return types.ActionContinue
}

func (r *responseContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	data, _, err := proxywasm.GetSharedData(sharedDataKey)
	if err != nil {
		return 0
	}

	buf := make([]byte, 8)
	ret := binary.LittleEndian.Uint64(data)
	binary.LittleEndian.PutUint64(buf, ret)

	proxywasm.LogInfof("返回体中获取shareData: %d", ret)

	body, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		return 0
	}
	if err != nil {
		proxywasm.LogErrorf("failed to get response body: %v", err)
		return types.ActionContinue
	}
	bodyStr := string(body)
	proxywasm.LogInfof("response body: %s", bodyStr)
	proxywasm.LogInfof("response body 是否结束: %t", endOfStream)

	return types.ActionContinue
}
