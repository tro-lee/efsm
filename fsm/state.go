package fsm

import (
	"context"
)

type State struct {
	Source      []string
	Destination string

	onEnterFunc func(*State)

	parallel bool
	fromAny  bool
	ctx      context.Context
	cancel   context.CancelFunc
}


func (st *State) To(dn string) *State {
	st.Destination = dn
	return st
}

func (st *State) FromAny() *State {
	st.fromAny = true
	return st
}

func (st *State) From(src ...string) *State {
	st.Source = src
	return st
}

func (st *State) FromStart() *State {
	st.Source = []string{"start"}
	return st
}

func (st *State) OnEnter(f func(s *State)) *State {
	st.onEnterFunc = f
	return st
}

func (st *State) Parallel(p bool) *State {
	st.parallel = p
	return st
}

func (st *State) Context() context.Context {
	if st.ctx != nil {
		return st.ctx
	}

	return context.Background()
}
