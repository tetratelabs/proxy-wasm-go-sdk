package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"regexp"
	"strings"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

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
	proxywasm.LogInfo("<---- pluginCx NewHttpContext ---->")
	return &responseContext{contextID: contextID}
}

type responseContext struct {
	contextID uint32
	types.DefaultHttpContext
}

func (r *responseContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfof("endOfStream: %t , body大小: %d ", endOfStream, bodySize)
	body, err := proxywasm.GetHttpResponseBody(0, bodySize)
	if err != nil {
		return 0
	}
	if err != nil {
		proxywasm.LogErrorf("failed to get response body: %v", err)
		return types.ActionContinue
	}
	bodyStr := string(body)
	proxywasm.LogInfof("original response body: %s", bodyStr)

	phonePat := ".*([^0-9]{1})(13|14|15|17|18|19)(\\d{9})([^0-9]{1}).*"
	replacePat := ".*(\\d{3})(\\d{4})(\\d{4}).*"
	phoneRegex := regexp.MustCompile(phonePat)
	replaceRegex := regexp.MustCompile(replacePat)

	for {
		subMatch := phoneRegex.FindStringSubmatch(bodyStr)
		if len(subMatch) != 0 {
			phoneNumber := subMatch[2] + subMatch[3]
			allString := replaceRegex.ReplaceAllString(phoneNumber, "$1****$3")
			bodyStr = strings.ReplaceAll(bodyStr, phoneNumber, allString)
		} else {
			break
		}
	}
	proxywasm.ReplaceHttpResponseBody([]byte(bodyStr))
	return types.ActionContinue
}
