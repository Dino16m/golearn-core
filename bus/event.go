package bus

type BaseEvent struct {
	Name    string
	Payload any
}

func (event BaseEvent) Data() any {
	return event.Payload
}

func (event BaseEvent) ID() string {
	return event.Name
}
