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

func (s *StateMachine) Exists() bool {
	return s.CurrentState != nil
}

func (s *StateMachine) IsValidStateChange(name string) (*State, error) {
	var st *State
	for _, state := range s.States {
		if state.Destination == name {
			st = state
		}
	}

	if st == nil {
		return st, fmt.Errorf("不存在目标状态: %v", name)
	}

	if st.fromAny || s.CurrentState == nil {
		return st, nil
	}

	for _, source := range st.Source {
		if source == s.CurrentState.Destination {
			return st, nil
		}
	}

	return st, fmt.Errorf("无法实现转移: %v > %v", s.CurrentState.Destination, st.Destination)
}

func (s *StateMachine) Transition(to string, data interface{}) (err error) {
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
		Data: data,
	}

	// 直接取消当前状态
	if s.CurrentState != nil && s.CurrentState.cancel != nil {
		s.CurrentState.cancel()
	}

	if s.ctx != nil && s.ctx.Err() != nil {
		return
	}
	s.transitions <- tr
	s.CurrentState = state
	return
}

// 创建状态
func (s *StateMachine) NewState() *State {
	st := &State{}
	s.States = append(s.States, st)

	return st
}

// 设置上下文
func (s *StateMachine) WithContext(ctx context.Context) *StateMachine {
	s.ctx, s.cancel = context.WithCancel(ctx)
	return s
}

func New() *StateMachine {
	return &StateMachine{
		transitions: make(chan *Transition, 1),
	}
}

// 启动状态机
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
