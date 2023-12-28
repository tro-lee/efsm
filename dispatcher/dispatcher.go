package dispatcher

import (
	"context"
	"log"
	"reflect"
)

type EventDispathcer struct {
	eventQueue  []Event
	handlerPool map[EventType][]EventHandler

	ctx context.Context

	initHandlers chan EventHandler
}

func New(ctx context.Context) *EventDispathcer {
	return &EventDispathcer{
		ctx:          ctx,
		initHandlers: make(chan EventHandler, 1),
		handlerPool:  make(map[EventType][]EventHandler),
		eventQueue:   make([]Event, 0),
	}
}

func (ed *EventDispathcer) Start() {
	log.Printf("\033[32mSuccess start dispatcher!\033[0m\n")
	ed.ctx = context.WithValue(ed.ctx, "dispather", ed)
	go func() {
		select {
		case <-ed.ctx.Done():
			return
		case handler := <-ed.initHandlers:
			log.Printf("\033[33mSuccess register Handler!: %v\033[0m\n", reflect.TypeOf(handler))
			handler.Start(ed.ctx)
			ed.handlerPool[handler.Type()] = append(ed.handlerPool[handler.Type()], handler)
		}
	}()
}

func (ed *EventDispathcer) AddEvent(event Event) {
	log.Printf("Add event: %v\n", event.Type())
	ed.eventQueue = append(ed.eventQueue, event)
}

func (ed *EventDispathcer) RegisterHandler(eh EventHandler) {
	log.Printf("Add handler: %v\n", reflect.TypeOf(eh))
	ed.initHandlers <- eh
}

func (ed *EventDispathcer) HandlerPool() map[EventType][]EventHandler {
	return ed.handlerPool
}

func (ed *EventDispathcer) EventQueue() []Event {
	return ed.eventQueue
}

func (ed *EventDispathcer) Pop() Event {
	result := ed.eventQueue[0]
	ed.eventQueue = ed.eventQueue[1:]
	return result
}
