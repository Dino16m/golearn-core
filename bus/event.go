package bus

type BaseEvent struct {
	Payload any
}

func (event BaseEvent) Data() any {
	return event.Payload
}
