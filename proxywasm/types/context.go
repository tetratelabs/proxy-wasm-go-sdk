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
	// In that case, must call Done() host call after the work is done to indicate that
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
	// NewTcpContext is used for creating TcpContext for each tcp streams.
	NewTcpContext(contextID uint32) TcpContext
	// NewHttpContext is used for creating HttpContext for each http streams.
	NewHttpContext(contextID uint32) HttpContext
}

// TcpContext corresponds to each TCP stream and is created by RootContext via NewTcpContext.
type TcpContext interface {
	// OnNewConnection is called when the tcp connection is established between Down and Upstreams.
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
	OnHttpRequestHeaders(numHeaders int, endOfStream bool) Action
	OnHttpRequestBody(bodySize int, endOfStream bool) Action
	OnHttpRequestTrailers(numTrailers int) Action
	OnHttpResponseHeaders(numHeaders int, endOfStream bool) Action
	OnHttpResponseBody(bodySize int, endOfStream bool) Action
	OnHttpResponseTrailers(numTrailers int) Action
	OnHttpStreamDone()
}

type (
	DefaultRootContext struct{}
	DefaultTcpContext  struct{}
	DefaultHttpContext struct{}
)

var (
	_ RootContext = &DefaultRootContext{}
	_ TcpContext  = &DefaultTcpContext{}
	_ HttpContext = &DefaultHttpContext{}
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
