package fsm

import (
	"context"
	"fmt"
)

type StateMachine struct {
	CurrentState *State
	States       []*State
	transitions  chan *Transition

	beforeFn func(*Transition)
	afterFn  func(*Transition)

	initialized bool

	ctx    context.Context
	cancel context.CancelFunc
}

func (s *StateMachine) Transitions() <-chan *Transition {
	return s.transitions
}

func (s *StateMachine) BeforeTransition(f func(*Transition)) {
	s.beforeFn = f
}

func (s *StateMachine) AfterTransition(f func(*Transition)) {
	s.afterFn = f
}

func (s *StateMachine) before(t *Transition) {
	if s.beforeFn != nil {
		s.beforeFn(t)
	}
}

func (s *StateMachine) after(t *Transition) {
	if s.afterFn != nil {
		s.afterFn(t)
	}
}

func (s *StateMachine) Find(st string) (state *State, err error) {
	for _, state := range s.States {
		if state.Destination == st {
			return state, nil
		}
	}

	return nil, fmt.Errorf("非法状态: %v", st)
}

func (s *StateMachine) Match(compare ...string) bool {
	if !s.Exists() {
		return false
	}

	for _, state := range compare {
		match := s.CurrentState.Destination == state
		if match {
			return true
		}
	}
	return false
}

func (s *StateMachine) Exists() bool {
	return s.CurrentState != nil
}

func (s *StateMachine) Start() {
	if s.initialized {
		return
	}

	s.initialized = true

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case t := <-s.transitions:
				if s.ctx.Err() != nil {
					return
				}

				s.before(t)
				t.do()
				s.after(t)
			}
		}
	}()
}

func (s *StateMachine) Name() string {
	if s.Exists() {
		return s.CurrentState.Destination
	}
	return ""
}

func (s *StateMachine) IsValidStateChange(name string) (*State, error) {
	st, err := s.Find(name)
	if err != nil {
		return st, err
	}

	if st.fromAny {
		return st, nil
	}

	if s.CurrentState == nil {
		return st, nil
	}

	for _, source := range st.Source {
		if source == s.CurrentState.Destination {
			return st, nil
		}
	}

	return st, fmt.Errorf("Invalid state change: %v > %v", s.CurrentState.Destination, st.Destination)
}

func (s *StateMachine) Transition(to string) (err error) {
	if s.Match(to) {
		return
	}

	state, err := s.IsValidStateChange(to)

	if err != nil {
		return
	}

	if s.ctx != nil {
		state.ctx, state.cancel = context.WithCancel(s.ctx)
	}

	tr := &Transition{
		From: s.CurrentState,
		To:   state,
	}

	if s.CurrentState != nil && s.CurrentState.cancel != nil {
		s.CurrentState.cancel()
	}

	if state.parallel {
		go tr.do()
	} else {
		if s.ctx != nil && s.ctx.Err() != nil {
			return
		}
		s.transitions <- tr
	}
	s.CurrentState = state
	return
}

func (s *StateMachine) NewState() *State {
	st := &State{}
	s.States = append(s.States, st)

	return st
}

func (s *StateMachine) WithContext(ctx context.Context) *StateMachine {
	s.ctx, s.cancel = context.WithCancel(ctx)
	return s
}

func New() *StateMachine {
	return &StateMachine{
		transitions: make(chan *Transition, 1),
	}
}
