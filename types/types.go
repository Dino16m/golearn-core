package types

import "github.com/dino16m/golearn-core/bus"

type UserCreated struct {
	bus.BaseEvent
}

const RefreshTokenKey = "refresh"
const AuthTokenKey = "auth"
