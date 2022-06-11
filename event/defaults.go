package event

import "github.com/dino16m/golearn-core/bus"

const (
	// expects map[string]interface{} with keys firstName, lastName, email,
	// and id which should have an int value
	UserCreated = EventName("usercreated")
)

type UserCreatedEvent struct {
	bus.BaseEvent
}

func NewUserCreatedEvent(payload any) UserCreatedEvent {
	return UserCreatedEvent{
		bus.BaseEvent{
			Payload: payload,
			Name:    "usercreated",
		},
	}
}
