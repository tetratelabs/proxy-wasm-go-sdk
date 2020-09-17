// Copyright 2020 Tetrate
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

package proxywasm

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type Context interface {
	OnDone() bool
	OnHttpCallResponse(numHeaders, bodySize, numTrailers int)
	OnLog()
}

type RootContext interface {
	Context
	OnConfigure(pluginConfigurationSize int) bool
	OnQueueReady(queueID uint32)
	OnTick()
	OnVMStart(vmConfigurationSize int) bool
}

type StreamContext interface {
	Context
	OnDownstreamData(dataSize int, endOfStream bool) types.Action
	OnDownstreamClose(peerType types.PeerType)
	OnNewConnection() types.Action
	OnUpstreamData(dataSize int, endOfStream bool) types.Action
	OnUpstreamClose(peerType types.PeerType)
}

type HttpContext interface {
	Context
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action
	OnHttpRequestBody(bodySize int, endOfStream bool) types.Action
	OnHttpRequestTrailers(numTrailers int) types.Action
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action
	OnHttpResponseBody(bodySize int, endOfStream bool) types.Action
	OnHttpResponseTrailers(numTrailers int) types.Action
}

type DefaultContext struct{}

var (
	_ Context       = DefaultContext{}
	_ RootContext   = DefaultContext{}
	_ StreamContext = DefaultContext{}
	_ HttpContext   = DefaultContext{}
)

// impl Context
func (d DefaultContext) OnDone() bool                     { return true }
func (d DefaultContext) OnHttpCallResponse(int, int, int) {}
func (d DefaultContext) OnLog()                           {}

// impl RootContext
func (d DefaultContext) OnConfigure(int) bool { return true }
func (d DefaultContext) OnQueueReady(uint32)  {}
func (d DefaultContext) OnTick()              {}
func (d DefaultContext) OnVMStart(int) bool   { return true }

// impl StreamContext
func (d DefaultContext) OnDownstreamData(int, bool) types.Action { return types.ActionContinue }
func (d DefaultContext) OnDownstreamClose(types.PeerType)        {}
func (d DefaultContext) OnNewConnection() types.Action           { return types.ActionContinue }
func (d DefaultContext) OnUpstreamData(int, bool) types.Action   { return types.ActionContinue }
func (d DefaultContext) OnUpstreamClose(types.PeerType)          {}

// impl HttpContext
func (d DefaultContext) OnHttpRequestHeaders(int, bool) types.Action  { return types.ActionContinue }
func (d DefaultContext) OnHttpRequestBody(int, bool) types.Action     { return types.ActionContinue }
func (d DefaultContext) OnHttpRequestTrailers(int) types.Action       { return types.ActionContinue }
func (d DefaultContext) OnHttpResponseHeaders(int, bool) types.Action { return types.ActionContinue }
func (d DefaultContext) OnHttpResponseBody(int, bool) types.Action    { return types.ActionContinue }
func (d DefaultContext) OnHttpResponseTrailers(int) types.Action      { return types.ActionContinue }
