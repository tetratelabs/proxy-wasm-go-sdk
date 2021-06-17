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

package types

type RootContext interface {
	OnQueueReady(queueID uint32)
	OnTick()
	OnVMStart(vmConfigurationSize int) OnVMStartStatus
	OnPluginStart(pluginConfigurationSize int) OnPluginStartStatus
	OnPluginDone() bool
	OnLog()

	// Child context factories
	NewStreamContext(contextID uint32) StreamContext
	NewHttpContext(contextID uint32) HttpContext
}

type StreamContext interface {
	OnDownstreamData(dataSize int, endOfStream bool) Action
	OnDownstreamClose(peerType PeerType)
	OnNewConnection() Action
	OnUpstreamData(dataSize int, endOfStream bool) Action
	OnUpstreamClose(peerType PeerType)
	OnStreamDone()
	OnLog()
}

type HttpContext interface {
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) Action
	OnHttpRequestBody(bodySize int, endOfStream bool) Action
	OnHttpRequestTrailers(numTrailers int) Action
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) Action
	OnHttpResponseBody(bodySize int, endOfStream bool) Action
	OnHttpResponseTrailers(numTrailers int) Action
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
func (*DefaultRootContext) OnQueueReady(uint32)           {}
func (*DefaultRootContext) OnTick()                       {}
func (*DefaultRootContext) OnVMStart(int) OnVMStartStatus { return OnVMStartStatusOK }
func (*DefaultRootContext) OnPluginStart(int) OnPluginStartStatus {
	return OnPluginStartStatusOK
}
func (*DefaultRootContext) OnPluginDone() bool                    { return true }
func (*DefaultRootContext) OnLog()                                {}
func (*DefaultRootContext) NewStreamContext(uint32) StreamContext { return nil }
func (*DefaultRootContext) NewHttpContext(uint32) HttpContext     { return nil }

// impl StreamContext
func (*DefaultStreamContext) OnDownstreamData(int, bool) Action { return ActionContinue }
func (*DefaultStreamContext) OnDownstreamClose(PeerType)        {}
func (*DefaultStreamContext) OnNewConnection() Action           { return ActionContinue }
func (*DefaultStreamContext) OnUpstreamData(int, bool) Action   { return ActionContinue }
func (*DefaultStreamContext) OnUpstreamClose(PeerType)          {}
func (*DefaultStreamContext) OnStreamDone()                     {}
func (*DefaultStreamContext) OnLog()                            {}

// impl HttpContext
func (*DefaultHttpContext) OnHttpRequestHeaders(int, bool) Action  { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestBody(int, bool) Action     { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestTrailers(int) Action       { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseHeaders(int, bool) Action { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseBody(int, bool) Action    { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseTrailers(int) Action      { return ActionContinue }
func (*DefaultHttpContext) OnHttpStreamDone()                      {}
func (*DefaultHttpContext) OnLog()                                 {}
