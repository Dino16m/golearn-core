package event

type UserCreated struct {
	Payload any
}

func NewUserCreatedEvent(payload any) UserCreated {
	return UserCreated{Payload: payload}
}
