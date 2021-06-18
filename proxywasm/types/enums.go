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

type Action uint32

const (
	ActionContinue Action = 0
	ActionPause    Action = 1
)

type PeerType uint32

const (
	PeerTypeUnknown PeerType = 0
	PeerTypeLocal   PeerType = 1
	PeerTypeRemote  PeerType = 2
)

type OnPluginStartStatus bool

const (
	OnPluginStartStatusOK     OnPluginStartStatus = true
	OnPluginStartStatusFailed OnPluginStartStatus = false
)

type OnVMStartStatus bool

const (
	OnVMStartStatusOK     OnVMStartStatus = true
	OnVMStartStatusFailed OnVMStartStatus = false
)
