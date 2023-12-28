package dispatcher

type Event interface {
	Type() EventType
	Data() []byte
}

type EventType string
