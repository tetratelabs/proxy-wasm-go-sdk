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

var currentState = &state{
	rootContexts:   make(map[uint32]RootContext),
	httpContexts:   make(map[uint32]HttpContext),
	streamContexts: make(map[uint32]StreamContext),
	callOuts:       make(map[uint32]uint32),
}

type state struct {
	newRootContext   func(contextID uint32) RootContext
	newStreamContext func(contextID uint32) StreamContext
	newHttpContext   func(contextID uint32) HttpContext
	rootContexts     map[uint32]RootContext
	httpContexts     map[uint32]HttpContext
	streamContexts   map[uint32]StreamContext
	activeContextID  uint32
	callOuts         map[uint32]uint32
}

func SetNewRootContext(f func(contextID uint32) RootContext) {
	currentState.newRootContext = f
}

func SetNewHttpContext(f func(contextID uint32) HttpContext) {
	currentState.newHttpContext = f
}

func SetNewStreamContext(f func(contextID uint32) StreamContext) {
	currentState.newStreamContext = f
}

func (s *state) createRootContext(contextID uint32) {
	var ctx RootContext
	if s.newRootContext == nil {
		ctx = &DefaultContext{}
	} else {
		ctx = s.newRootContext(contextID)
	}

	s.rootContexts[contextID] = ctx
}

func (s *state) createStreamContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.streamContexts[contextID]; ok {
		panic("context id duplicated")
	}

	s.streamContexts[contextID] = s.newStreamContext(contextID)
}

func (s *state) createHttpContext(contextID uint32, rootContextID uint32) {
	if _, ok := s.rootContexts[rootContextID]; !ok {
		panic("invalid root context id")
	}

	if _, ok := s.httpContexts[contextID]; ok {
		panic("context id duplicated")
	}

	s.httpContexts[contextID] = s.newHttpContext(contextID)
}

func (s *state) registerCallout(calloutID uint32) {
	if _, ok := s.callOuts[calloutID]; ok {
		panic("duplicated calloutID")
	}

	s.callOuts[calloutID] = s.activeContextID
}

func (s *state) setActiveContextID(contextID uint32) {
	// TODO: should we do this inline (possibly for performance)?
	s.activeContextID = contextID
}
