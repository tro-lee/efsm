package main

import (
	"context"
	"learn/efsm/dispatcher"
	"learn/efsm/handlers"
	"time"
)

func main() {
	ctx := context.Background()

	globalDispatcher := dispatcher.New(ctx)
	globalDispatcher.Start()

	globalDispatcher.RegisterHandler(handlers.NewAllocation())
	globalDispatcher.RegisterHandler(handlers.NewParser())

	time.Sleep(1 * time.Second)
	globalDispatcher.AddEvent(&handlers.ParserEvent{})

	globalDispatcher.AddEvent(&handlers.ParserEvent{})

	for {
	}
}
