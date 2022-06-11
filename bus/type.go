package bus

type Event interface {
	ID() string
	Data() any
}

type Listener interface {
	Handle(event Event)
}
