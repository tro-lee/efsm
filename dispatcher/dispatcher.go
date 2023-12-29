package dispatcher

import (
	"context"
	"learn/efsm/out"
	"log"
	"reflect"
)

type EventDispathcer struct {
	eventQueue  chan Event
	handlerPool map[EventType][]EventHandler

	ctx context.Context

	initHandlers chan EventHandler
}

func New(ctx context.Context) *EventDispathcer {
	return &EventDispathcer{
		ctx:          ctx,
		initHandlers: make(chan EventHandler, 20),
		handlerPool:  make(map[EventType][]EventHandler),
		eventQueue:   make(chan Event, 50),
	}
}

func (ed *EventDispathcer) Start() {
	log.Printf("\033[32mSuccess start dispatcher!\033[0m\n")
	ed.ctx = context.WithValue(ed.ctx, "dispather", ed)
	go func() {
		for {
			select {
			case <-ed.ctx.Done():
				return
			case handler := <-ed.initHandlers:
				out.Success("Success register Handler!: %v", reflect.TypeOf(handler))
				handler.Start(ed.ctx)
				ed.handlerPool[handler.Type()] = append(ed.handlerPool[handler.Type()], handler)
			}
		}
	}()
}

func (ed *EventDispathcer) AddEvent(event Event) {
	ed.eventQueue <- event
}

func (ed *EventDispathcer) RegisterHandler(eh EventHandler) {
	ed.initHandlers <- eh
}

func (ed *EventDispathcer) HandlerPool() map[EventType][]EventHandler {
	return ed.handlerPool
}

func (ed *EventDispathcer) EventQueue() chan Event {
	return ed.eventQueue
}
