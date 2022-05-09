package main

import (
	"fmt"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
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
	configuration pluginConfiguration
}

type pluginConfiguration struct {
	// desensitizeTypes support values includes honeNumber/idCard
	desensitizeTypes []string
}

// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &responseContext{contextID: contextID, desensitizeTypes: ctx.configuration.desensitizeTypes}
}

type responseContext struct {
	contextID uint32
	types.DefaultHttpContext
	desensitizeTypes []string
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
	desensitizeTypes := jsonData.Get("desensitizeTypes").Array()
	for _, desensitizeType := range desensitizeTypes {
		config.desensitizeTypes = append(config.desensitizeTypes, desensitizeType.Str)
	}

	return *config, nil
}

func (r *responseContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {

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

	enablePhoneNumber := false
	enableIdCard := false

	if len(r.desensitizeTypes) != 0 {
		enablePhoneNumber = isContain(r.desensitizeTypes, "PhoneNumber")
		enableIdCard = isContain(r.desensitizeTypes, "IdCard")
	}

	if enablePhoneNumber {
		bodyStr = PhoneNumberDesensitize(bodyStr)
	}

	if enableIdCard {
		bodyStr = IdCardDesensitize(bodyStr)
	}

	proxywasm.ReplaceHttpResponseBody([]byte(bodyStr))
	return types.ActionContinue
}

func isContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func postCodeDesensitize(body string) string {
	postCodePat := "^[0-9]\\d{5}$"
	replacePat := "(\\d{2})(\\d{2})(\\d{2})"
	postCodeRegex := regexp.MustCompile(postCodePat)
	replaceRegex := regexp.MustCompile(replacePat)
	for {
		subMatch := postCodeRegex.FindStringSubmatch(body)
		if len(subMatch) != 0 {
			allString := replaceRegex.ReplaceAllString(subMatch[0], "**$2$3")
			body = strings.ReplaceAll(body, subMatch[0], allString)
		} else {
			break
		}
	}
	return body
}

func IdCardDesensitize(body string) string {
	idCardPat := "(([1-6]\\d{5})(19\\d{2}|20\\d{2})(0[1-9]|1[012])(0[1-9]|[12]\\d|3[01])(\\d{3}[\\dxX]))"
	replacePat := "(\\d{6})(\\d{8})(\\w{4})"
	idCardRegex := regexp.MustCompile(idCardPat)
	replaceRegex := regexp.MustCompile(replacePat)
	for {
		subMatch := idCardRegex.FindStringSubmatch(body)
		if len(subMatch) != 0 {
			allString := replaceRegex.ReplaceAllString(subMatch[0], "$1********$3")
			body = strings.ReplaceAll(body, subMatch[0], allString)
		} else {
			break
		}
	}
	return body
}

func PhoneNumberDesensitize(body string) string {
	phonePat := ".*([^0-9]{1})(13|14|15|17|18|19)(\\d{9})([^0-9]{1}).*"
	replacePat := ".*(\\d{3})(\\d{4})(\\d{4}).*"
	phoneRegex := regexp.MustCompile(phonePat)
	replaceRegex := regexp.MustCompile(replacePat)
	for {
		subMatch := phoneRegex.FindStringSubmatch(body)
		if len(subMatch) != 0 {
			phoneNumber := subMatch[2] + subMatch[3]
			allString := replaceRegex.ReplaceAllString(phoneNumber, "$1****$3")
			body = strings.ReplaceAll(body, phoneNumber, allString)
		} else {
			break
		}
	}
	return body
}
