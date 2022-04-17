package config

import "time"

// JwtOptions ...
type JwtOptions struct {
	Realm       string
	Key         string
	Timeout     time.Duration
	MaxRefresh  time.Duration
	IdentityKey string
}

type AuthConfig struct {
	UserIdClaim        string
	AuthUserContextKey string
}

var UserIdClaim string
var AuthUserContextKey string

func init() {
	UserIdClaim = "uid"
	AuthUserContextKey = "authusercontext"
}

func Setup(cfg AuthConfig) {
	UserIdClaim = cfg.UserIdClaim
	AuthUserContextKey = cfg.AuthUserContextKey
}
