package types

import (
	"github.com/gin-gonic/gin"
)

// AuthUserManager ...
type AuthUserManager interface {
	GetAuthUser(c *gin.Context) (AuthUser, error)
}

// AuthUser ...
type AuthUser interface {
	GetPassword() string
	GetId() int
	EmailVerified()
	SetPassword(string)
	// GetName returns user name in the format FirstName LastName
	GetName() (string, string)
	GetEmail() string
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}
