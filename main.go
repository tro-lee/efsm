package main

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/handlers"
)

func main() {
	ctx := context.Background()

	globalDispatcher := dispatcher.New(ctx)
	allocationHandler := handlers.NewAllocation()

	globalDispatcher.RegisterHandler(allocationHandler)
	globalDispatcher.Start()

	globalDispatcher.AddEvent(&EventOne{})
	for {
	}
}

type EventOne struct {
}

func (e EventOne) Type() dispatcher.EventType {
	return "EventOne"
}

func (e EventOne) Data() []byte {
	return []byte("EventOne")
}
