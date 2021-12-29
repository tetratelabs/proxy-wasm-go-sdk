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

import "errors"

// Action represents the action which Wasm contexts expects hosts to take.
type Action uint32

const (
	// ActionContinue means that the host continues the processing.
	ActionContinue Action = 0
	// ActionPause means that the host pauses the processing.
	ActionPause Action = 1
)

// PeerType represents the type of a peer of a connection.
type PeerType uint32

const (
	// PeerTypeUnknown means the type of a peer is unknwo
	PeerTypeUnknown PeerType = 0
	// PeerTypeLocal means the type of a peer is local (i.e. proxy)
	PeerTypeLocal PeerType = 1
	// PeerTypeRemote means the type of a peer is remote (i.e. remote client)
	PeerTypeRemote PeerType = 2
)

// OnVMStartStatus is the tyep of status returned by OnVMStart
type OnVMStartStatus bool

const (
	// OnVMStartStatusOK indicates that VMContext.OnVMStartStatus succeeded.
	OnVMStartStatusOK OnVMStartStatus = true
	// OnVMStartStatusFailed indicates that VMContext.OnVMStartStatus failed.
	// The further processing for this VM never happens, and hosts would
	// delete this VM.
	OnVMStartStatusFailed OnVMStartStatus = false
)

// OnPluginStartStatus is the tyep of status returned by OnPluginStart
type OnPluginStartStatus bool

const (
	// OnPluginStartStatusOK indicates that PluginContext.OnPluginStart succeeded.
	OnPluginStartStatusOK OnPluginStartStatus = true
	// OnPluginStartStatusFailed indicates that PluginContext.OnPluginStart failed.
	// The further processing for that plugin context never happens.
	OnPluginStartStatusFailed OnPluginStartStatus = false
)

var (
	// ErrorStatusNotFound means not found for various hostcalls.
	ErrorStatusNotFound = errors.New("error status returned by host: not found")
	// ErrorStatusNotFound means the arguments for a hostcall are invalid.
	ErrorStatusBadArgument = errors.New("error status returned by host: bad argument")
	// ErrorStatusNotFound means the target queue of DequeueSharedQueue call is empty.
	ErrorStatusEmpty = errors.New("error status returned by host: empty")
	// ErrorStatusNotFound means a given CAS value for SetSharedData is mismatched
	// with the current value. That indicates that other Wasm VMs has already succeeded
	// to set a value on the same key and the current CAS for the key is incremented.
	// Having retry logic in the face of this error is recommended.
	ErrorStatusCasMismatch = errors.New("error status returned by host: cas mismatch")
	// ErrorStatusNotFound indicates an internal falure in hosts.
	// In the face of this error, there's nothing we could do in the Wasm VM.
	// Recommend simply abort or panic then.
	ErrorInternalFailure = errors.New("error status returned by host: internal failure")
	// ErrorUnimplemented indicates the API is not implemented in the host yet.
	ErrorUnimplemented = errors.New("error status returned by host: unimplemented")
)
