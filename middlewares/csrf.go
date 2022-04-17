package middlewares

import (
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

// CSRFMiddleware ...
type CSRFMiddleware struct {
	key    string
	secure bool
	opts   sessions.Options
}

// NewCSRFMiddleware construct the middleware
func NewCSRFMiddleware(key string,
	env string, opts sessions.Options) CSRFMiddleware {
	secure := strings.ToLower(env) == "production"

	return CSRFMiddleware{
		key:    key,
		secure: secure,
		opts:   opts,
	}
}

// Guard returns the main csrf middleware
func (m CSRFMiddleware) Guard() gin.HandlerFunc {
	mw := csrf.Middleware(csrf.Options{
		Secret: m.key,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	})
	return mw
}

// TokenInjector is a terminable middleware used to inject an csrftoken cookie
// in all responses sent by the app
func (m CSRFMiddleware) TokenInjector(c *gin.Context) {
	csrfToken := csrf.GetToken(c)
	c.SetCookie(
		"csrftoken",
		csrfToken,
		m.opts.MaxAge,
		m.opts.Path,
		m.opts.Domain,
		m.opts.Secure,
		m.opts.HttpOnly,
	)
	c.Next()
}
