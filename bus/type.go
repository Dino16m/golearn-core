package bus

type Event interface {
}

type Listener interface {
	Handle(event Event)
}
