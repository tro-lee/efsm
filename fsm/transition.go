package fsm

type Transition struct {
	From *State
	To   *State
	Data interface{}
}

func (t *Transition) do() {
	if t.To.onEnterFunc != nil {
		t.To.onEnterFunc(t.From, t.Data)
	}
}
