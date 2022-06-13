package bus

type Event interface {
	Data() any
}

type Listener interface {
	Handle(event Event)
}
