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

// There are three types of "contexts" you are supposed to implement for writing Proxy-Wasm plugins with this SDK.
// They are called RootContext, TcpContext and HttpContext, and their relationship can be described as the following diagram:
//
//              ╱ TcpContext = handling each Tcp stream
//             ╱
//            ╱ 1: N
//  RootContext
//            ╲ 1: N
//             ╲
//              ╲ Http = handling each Http stream
//
// In other words, RootContex is the parent of others, and responsible for creating Tcp and Http contexts
// corresponding to each streams if it is configured for running as a Http/Tcp stream plugin.
// Given that, RootContext is the primary interface everyone has to implement.
//
// RootContext corresponds to each different plugin configurations (config.configuration),
// and the root context created first is specially treated as a VM context, which can handle
// vm_config.configuration during OnVMStart call to do VM-wise initialization.
type RootContext interface {
	// OnVMStart is called after the VM is created and main function is called.
	// During this call, GetVmConfiguration hostcall is available and can be used to
	// retrieve the configuration set at vm_config.configuration.
	//
	// Note that **only one root cnotext is called on this function**.
	// That is because there's Wasm VM: RootContext = 1: N correspondence, and
	// the firstly created root context of these root contexts will be treated
	// as a *VM context* on which OnVMStart is invoked by host.
	// In other words, vm_config.configuration is only available for only one root context
	// which is created first.
	OnVMStart(vmConfigurationSize int) OnVMStartStatus

	// OnPluginStart is called on all root contexts (after OnVmStart if this is the VM context).
	// During this call, hostcalls.getPluginConfiguration is available and can be used to
	// retrieve the configuration set at config.configuration in envoy.yaml
	OnPluginStart(pluginConfigurationSize int) OnPluginStartStatus

	// onPluginDone is called right before root contexts are deleted by hosts.
	// Return false to indicate it's in a pending state to do some more work left.
	// In that case, must call PluginDone() host call after the work is done to indicate that
	// hosts can kill this contexts.
	OnPluginDone() bool

	// OnQueueReady is called when the queue is ready after calling RegisterQueue hostcall.
	// Note that the queue is dequeued by another VM running in another thread, so possibly
	// the queue is empty during OnQueueReady even if it is not dequeued by this VM.
	OnQueueReady(queueID uint32)

	// OnTick is called when SetTickPeriodMilliSeconds hostcall is called by this root context.
	// This can be used for doing some asynchronous tasks in parallel to stream processing.
	OnTick()

	// The following functions are used for creating contexts on streams,
	// and developers *must* implement either of them corresponding to
	// extension points. For example, if you configure this root context is running
	// at Http filters, then NewHttpContext must be implemented. Same goes for
	// Tcp filters.
	//
	// NewTcpContext is used for creating TcpContext for each Tcp streams.
	NewTcpContext(contextID uint32) TcpContext
	// NewHttpContext is used for creating HttpContext for each Http streams.
	NewHttpContext(contextID uint32) HttpContext
}

// TcpContext corresponds to each Tcp stream and is created by RootContext via NewTcpContext.
type TcpContext interface {
	// OnNewConnection is called when the Tcp connection is established between Down and Upstreams.
	OnNewConnection() Action

	// OnDownstreamData is called when a data frame arrives from the downstream connection.
	OnDownstreamData(dataSize int, endOfStream bool) Action

	// OnDownstreamClose is called when the downstream connection is closed.
	OnDownstreamClose(peerType PeerType)

	/// OnUpstreamData is called when a data frame arrives from the upstream connection.
	OnUpstreamData(dataSize int, endOfStream bool) Action

	// OnUpstreamClose is called when the upstream connection is closed.
	OnUpstreamClose(peerType PeerType)

	// OnStreamDone is called before the host deletes this context.
	// You can retreive the stream information (such as remote addesses, etc.) during this calls
	// This can be used for implementing logging feature.
	OnStreamDone()
}

// HttpContext corresponds to each Http stream and is created by RootContext via NewHttpContext.
type HttpContext interface {
	// OnHttpRequestHeaders is called when request headers arrives.
	// Return types.ActionPause if you want to stop sending headers to upstream.
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) Action

	// OnHttpRequestBody is called when a request body *frame* arrives.
	// Note that this is possibly called multiple times until we see end_of_stream = true.
	// Return types.ActionPause if you want to buffer the body and stop sending body to upstream.
	// Even after returning types.ActionPause, this will be called when a unseen frame arrives.
	OnHttpRequestBody(bodySize int, endOfStream bool) Action

	// OnHttpRequestTrailers is called when request trailers arrives.
	// Return types.ActionPause if you want to stop sending trailers to upstream.
	OnHttpRequestTrailers(numTrailers int) Action

	// OnHttpResponseHeaders is called when response headers arrives.
	// Return types.ActionPause if you want to stop sending headers to downstream.
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) Action

	// OnHttpResponseBody is called when a response body *frame* arrives.
	// Note that this is possibly called multiple times until we see end_of_stream = true.
	// Return types.ActionPause if you want to buffer the body and stop sending body to downtream.
	// Even after returning types.ActionPause, this will be called when a unseen frame arrives.
	OnHttpResponseBody(bodySize int, endOfStream bool) Action

	// OnHttpResponseTrailers is called when response trailers arrives.
	// Return types.ActionPause if you want to stop sending trailers to downstream.
	OnHttpResponseTrailers(numTrailers int) Action

	// OnHttpStreamDone is called before the host deletes this context.
	// You can retreive the HTTP request/response information (such headers, etc.) during this calls.
	// This can be used for implementing logging feature.
	OnHttpStreamDone()
}

// DefaultContexts are no-op implementation of contexts.
// Users can embed them into their custom contexts, so that
// they only have to implement methods they want.
type (
	// DefaultRootContext provides the no-op implementation of RootContext interface.
	DefaultRootContext struct{}

	// DefaultTcpContext provides the no-op implementation of TcpContext interface.
	DefaultTcpContext struct{}

	// DefaultHttpContext provides the no-op implementation of HttpContext interface.
	DefaultHttpContext struct{}
)

// impl RootContext
func (*DefaultRootContext) OnQueueReady(uint32)           {}
func (*DefaultRootContext) OnTick()                       {}
func (*DefaultRootContext) OnVMStart(int) OnVMStartStatus { return OnVMStartStatusOK }
func (*DefaultRootContext) OnPluginStart(int) OnPluginStartStatus {
	return OnPluginStartStatusOK
}
func (*DefaultRootContext) OnPluginDone() bool                { return true }
func (*DefaultRootContext) NewTcpContext(uint32) TcpContext   { return nil }
func (*DefaultRootContext) NewHttpContext(uint32) HttpContext { return nil }

// impl TcpContext
func (*DefaultTcpContext) OnDownstreamData(int, bool) Action { return ActionContinue }
func (*DefaultTcpContext) OnDownstreamClose(PeerType)        {}
func (*DefaultTcpContext) OnNewConnection() Action           { return ActionContinue }
func (*DefaultTcpContext) OnUpstreamData(int, bool) Action   { return ActionContinue }
func (*DefaultTcpContext) OnUpstreamClose(PeerType)          {}
func (*DefaultTcpContext) OnStreamDone()                     {}

// impl HttpContext
func (*DefaultHttpContext) OnHttpRequestHeaders(int, bool) Action  { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestBody(int, bool) Action     { return ActionContinue }
func (*DefaultHttpContext) OnHttpRequestTrailers(int) Action       { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseHeaders(int, bool) Action { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseBody(int, bool) Action    { return ActionContinue }
func (*DefaultHttpContext) OnHttpResponseTrailers(int) Action      { return ActionContinue }
func (*DefaultHttpContext) OnHttpStreamDone()                      {}

var (
	_ RootContext = &DefaultRootContext{}
	_ TcpContext  = &DefaultTcpContext{}
	_ HttpContext = &DefaultHttpContext{}
)
