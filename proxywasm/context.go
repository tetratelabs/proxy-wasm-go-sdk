// Copyright 2020-2021 Tetrate
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

type RootContext interface {
	OnQueueReady(queueID uint32)
	OnTick()
	OnVMStart(vmConfigurationSize int) types.OnVMStartStatus
	OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus
	OnPluginDone() bool
	OnLog()

	// Child context factories
	NewStreamContext(contextID uint32) StreamContext
	NewHttpContext(contextID uint32) HttpContext
}

type StreamContext interface {
	OnDownstreamData(dataSize int, endOfStream bool) types.Action
	OnDownstreamClose(peerType types.PeerType)
	OnNewConnection() types.Action
	OnUpstreamData(dataSize int, endOfStream bool) types.Action
	OnUpstreamClose(peerType types.PeerType)
	OnStreamDone()
	OnLog()
}

type HttpContext interface {
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action
	OnHttpRequestBody(bodySize int, endOfStream bool) types.Action
	OnHttpRequestTrailers(numTrailers int) types.Action
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action
	OnHttpResponseBody(bodySize int, endOfStream bool) types.Action
	OnHttpResponseTrailers(numTrailers int) types.Action
	OnHttpStreamDone()
	OnLog()
}

type (
	DefaultRootContext   struct{}
	DefaultStreamContext struct{}
	DefaultHttpContext   struct{}
)

var (
	_ RootContext   = &DefaultRootContext{}
	_ StreamContext = &DefaultStreamContext{}
	_ HttpContext   = &DefaultHttpContext{}
)

// impl RootContext
func (*DefaultRootContext) OnQueueReady(uint32)                 {}
func (*DefaultRootContext) OnTick()                             {}
func (*DefaultRootContext) OnVMStart(int) types.OnVMStartStatus { return types.OnVMStartStatusOK }
func (*DefaultRootContext) OnPluginStart(int) types.OnPluginStartStatus {
	return types.OnPluginStartStatusOK
}
func (*DefaultRootContext) OnPluginDone() bool                    { return true }
func (*DefaultRootContext) OnLog()                                {}
func (*DefaultRootContext) NewStreamContext(uint32) StreamContext { return nil }
func (*DefaultRootContext) NewHttpContext(uint32) HttpContext     { return nil }

// impl StreamContext
func (*DefaultStreamContext) OnDownstreamData(int, bool) types.Action { return types.ActionContinue }
func (*DefaultStreamContext) OnDownstreamClose(types.PeerType)        {}
func (*DefaultStreamContext) OnNewConnection() types.Action           { return types.ActionContinue }
func (*DefaultStreamContext) OnUpstreamData(int, bool) types.Action   { return types.ActionContinue }
func (*DefaultStreamContext) OnUpstreamClose(types.PeerType)          {}
func (*DefaultStreamContext) OnStreamDone()                           {}
func (*DefaultStreamContext) OnLog()                                  {}

// impl HttpContext
func (*DefaultHttpContext) OnHttpRequestHeaders(int, bool) types.Action  { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpRequestBody(int, bool) types.Action     { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpRequestTrailers(int) types.Action       { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpResponseHeaders(int, bool) types.Action { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpResponseBody(int, bool) types.Action    { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpResponseTrailers(int) types.Action      { return types.ActionContinue }
func (*DefaultHttpContext) OnHttpStreamDone()                            {}
func (*DefaultHttpContext) OnLog()                                       {}
