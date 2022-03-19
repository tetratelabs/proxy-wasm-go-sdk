package main

import (
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	// SetVMContext is the entrypoint for setting up this entire Wasm VM.
	// Please make sure that this entrypoint be called during "main()" function, otherwise
	// this VM would fail.
	proxywasm.SetVMContext(&vmContext{})
}

// vmContext implements types.VMContext interface of proxy-wasm-go SDK.
type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// pluginContext implements types.PluginContext interface of proxy-wasm-go SDK.
type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	configuration pluginConfiguration
}

// pluginConfiguration is a type to represent an example configuration for this wasm plugin.
type pluginConfiguration struct {
	// Example configuration field.
	// The plugin will validate if those fields exist in the json payload.
	requiredKeys []string
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	config, err := parsePluginConfiguration(data)
	if err != nil {
		proxywasm.LogCriticalf("error parsing plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.configuration = config
	return types.OnPluginStartStatusOK
}

// parsePluginConfiguration parses the json plugin confiuration data and returns pluginConfiguration.
// Note that this parses the json data by gjson, since TinyGo doesn't support encoding/json.
// You can also try https://github.com/mailru/easyjson, which supports decoding to a struct.
func parsePluginConfiguration(data []byte) (pluginConfiguration, error) {
	if len(data) == 0 {
		return pluginConfiguration{}, nil
	}

	config := &pluginConfiguration{}
	if !gjson.ValidBytes(data) {
		return pluginConfiguration{}, fmt.Errorf("the plugin configuration is not a valid json: %q", string(data))
	}

	jsonData := gjson.ParseBytes(data)
	requiredKeys := jsonData.Get("requiredKeys").Array()
	for _, requiredKey := range requiredKeys {
		config.requiredKeys = append(config.requiredKeys, requiredKey.Str)
	}

	return *config, nil
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &payloadValidationContext{requiredKeys: ctx.configuration.requiredKeys}
}

// payloadValidationContext implements types.HttpContext interface of proxy-wasm-go SDK.
type payloadValidationContext struct {
	// Embed the default root http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	totalRequestBodySize int
	requiredKeys         []string
}

var _ types.HttpContext = (*payloadValidationContext)(nil)

// Override types.DefaultHttpContext.
func (*payloadValidationContext) OnHttpRequestHeaders(numHeaders int, _ bool) types.Action {
	contentType, err := proxywasm.GetHttpRequestHeader("content-type")
	if err != nil || contentType != "application/json" {
		// If the header doesn't have the expected content value, send the 403 response,
		if err := proxywasm.SendHttpResponse(403, nil, []byte("content-type must be provided"), -1); err != nil {
			panic(err)
		}
		// and terminates the further processing of this traffic by ActionPause.
		return types.ActionPause
	}

	// ActionContinue lets the host continue the processing the body.
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *payloadValidationContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	ctx.totalRequestBodySize += bodySize
	if !endOfStream {
		// OnHttpRequestBody may be called each time a part of the body is received.
		// Wait until we see the entire body to replace.
		return types.ActionPause
	}

	body, err := proxywasm.GetHttpRequestBody(0, ctx.totalRequestBodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get request body: %v", err)
		return types.ActionContinue
	}
	if !ctx.validatePayload(body) {
		// If the validation fails, send the 403 response,
		if err := proxywasm.SendHttpResponse(403, nil, []byte("invalid payload"), -1); err != nil {
			proxywasm.LogErrorf("failed to send the 403 response: %v", err)
		}
		// and terminates this traffic.
		return types.ActionPause
	}

	return types.ActionContinue
}

// validatePayload validates the given json payload.
// Note that this function parses the json data by gjson, since TinyGo doesn't support encoding/json.
func (ctx *payloadValidationContext) validatePayload(body []byte) bool {
	if !gjson.ValidBytes(body) {
		proxywasm.LogErrorf("body is not a valid json: %q", string(body))
		return false
	}
	jsonData := gjson.ParseBytes(body)

	// Do any validation on the json. Check if required keys exist here as an example.
	// The required keys are configurable via the plugin configuration.
	for _, requiredKey := range ctx.requiredKeys {
		if !jsonData.Get(requiredKey).Exists() {
			proxywasm.LogErrorf("required key (%v) is missing: %v", requiredKey, jsonData)
			return false
		}
	}

	return true
}
