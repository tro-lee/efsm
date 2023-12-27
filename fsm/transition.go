package fsm

type Transition struct {
	From *State
	To   *State
}

func (t *Transition) do() {
	if t.To.onEnterFunc != nil {
		t.To.onEnterFunc(t.To)
	}
}
