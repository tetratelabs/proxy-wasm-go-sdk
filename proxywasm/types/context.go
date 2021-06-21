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

// There are four types of these intefaces which you are supposed to implement in order to extend your network proxies.
// They are VMContext, PluginContext, TcpContext and HttpContext, and their relationship can be described as the following diagram:
//
//                          Wasm Virtual Machine(VM)
//                     (corresponds to VM configuration)
//  ┌────────────────────────────────────────────────────────────────────────────┐
//  │                                                      TcpContext            │
//  │                                                  ╱ (Each Tcp stream)       │
//  │                                                 ╱                          │
//  │                      1: N                      ╱ 1: N                      │
//  │       VMContext  ──────────  PluginContext                                 │
//  │  (VM configuration)     (Plugin configuration) ╲ 1: N                      │
//  │                                                 ╲                          │
//  │                                                  ╲   HttpContext           │
//  │                                                   (Each Http stream)       │
//  └────────────────────────────────────────────────────────────────────────────┘
//
// To summarize,
//
// 1) VMContext corresponds to each Wasm Virtual Machine, and only one VMContext exists in each VM.
// Note that in Envoy, Wasm VMs are created per "vm_config" field in envoy.yaml. For example having different "vm_config.configuration" fields
// results in multiple VMs being created and each of them corresponds to each "vm_config.configuration".
//
// 2) VMContext is parent of PluginContext, and is responsible for creating arbitrary number of PluginContexts.
//
// 3) PluginContext corresponds to each plugin configurations in the host. In Envoy, each plugin configuration is given at HttpFilter or NetworkFilter
// on listeners. That said, a plugin context corresponds to a Http or Network filter on a litener and is in charge of creating "filter instances" for
// each Http or Tcp streams. And these "filter instances" are HttpContexts or TcpContexts.
//
// 4) PluginContext is parent of TcpContext and HttpContexts, and is responsible for creating arbitrary number of these contexts.
//
// 5) TcpContext is responsible for handling each Tcp stream events.
//
// 6) HttpContext is responsible for handling each Http stream events.
//
//
// VMContext corresponds to each Wasm VM machine and its configuration. Thefore,
// this is the entrypoint for extending your network proxy.
// Its lifetime is exactly the same as Wasm Virtual Machines on the host.
type VMContext interface {
	// OnVMStart is called after the VM is created and main function is called.
	// During this call, GetVmConfiguration hostcall is available and can be used to
	// retrieve the configuration set at vm_config.configuration.
	// This is mainly used for doing Wasm VM-wise initialization.
	OnVMStart(vmConfigurationSize int) OnVMStartStatus

	// NewPluginContext is used for creating PluginContext for each plugin configurations.
	NewPluginContext(contextID uint32) PluginContext
}

// PluginContext corresponds to each different plugin configurations (config.configuration).
// Each configuration is usually given at each http/tcp filter in a listener in the hosts,
// so PluginContext is responsible for creating "filter instances" for each Tcp/Http streams on the listener.
type PluginContext interface {
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

// TcpContext corresponds to each Tcp stream and is created by PluginContext via NewTcpContext.
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

// HttpContext corresponds to each Http stream and is created by PluginContext via NewHttpContext.
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
	// DefaultPluginContext provides the no-op implementation of PluginContext interface.
	DefaultPluginContext struct{}

	// DefaultTcpContext provides the no-op implementation of TcpContext interface.
	DefaultTcpContext struct{}

	// DefaultHttpContext provides the no-op implementation of HttpContext interface.
	DefaultHttpContext struct{}
)

// impl PluginContext
func (*DefaultPluginContext) OnQueueReady(uint32) {}
func (*DefaultPluginContext) OnTick()             {}
func (*DefaultPluginContext) OnPluginStart(int) OnPluginStartStatus {
	return OnPluginStartStatusOK
}
func (*DefaultPluginContext) OnPluginDone() bool                { return true }
func (*DefaultPluginContext) NewTcpContext(uint32) TcpContext   { return nil }
func (*DefaultPluginContext) NewHttpContext(uint32) HttpContext { return nil }

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
	_ PluginContext = &DefaultPluginContext{}
	_ TcpContext    = &DefaultTcpContext{}
	_ HttpContext   = &DefaultHttpContext{}
)
