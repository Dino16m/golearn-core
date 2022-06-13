package event

import "github.com/dino16m/golearn-core/bus"

type UserCreated struct {
	bus.BaseEvent
}

func NewUserCreatedEvent(payload any) UserCreated {
	return UserCreated{
		bus.BaseEvent{
			Payload: payload,
		},
	}
}
