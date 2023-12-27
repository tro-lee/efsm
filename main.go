package main

import (
	"context"
	"learn/efsm/fsm"
)

func main() {
	ctx := context.Background()

	f := fsm.New().WithContext(ctx)

	f.NewState().FromStart().To("state1").OnEnter(func(s *fsm.State) {
	})
}
