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

// CORSConfig contains the settings used to setup CORS for the app
type CORSConfig struct {
	AllowAllOrigins bool
	AllowOrigins    []string

	// AllowMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowMethods []string

	// AllowHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposeHeaders []string

	// Allows to add origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
	AllowWildcard bool
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
