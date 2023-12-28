package dispatcher

import "context"

type EventHandler interface {
	Start(ctx context.Context)

	Status() EventHandlerStatus
	Type() EventType

	Handle(e Event)
}

type EventHandlerStatus string

const (
	Init    EventHandlerStatus = "Init"
	Ready   EventHandlerStatus = "Ready"
	Running EventHandlerStatus = "Running"
	Close   EventHandlerStatus = "Close"
	Error   EventHandlerStatus = "Error"
)
