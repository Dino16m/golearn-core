package middlewares

import (
	"strings"

	"github.com/dino16m/golearn-core/config"
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/services"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	FindAuthUser(username interface{}) (interface{}, errors.ApplicationError)
}

type JWTAuthMiddleware struct {
	authAdapter services.JWTAuthService
	userRepo    UserRepository
}

func NewJWTAuthMiddleware(authService services.JWTAuthService, userRepo UserRepository) JWTAuthMiddleware {
	return JWTAuthMiddleware{authAdapter: authService, userRepo: userRepo}
}

func (m JWTAuthMiddleware) Authorize(c *gin.Context) {
	authorization, ok := c.Request.Header["Authorization"]
	if !ok {
		errorResponse(c, errors.UnauthorizedError("Unauthorized"))
		return
	}
	header := authorization[len(authorization)-1]
	token := strings.Split(header, " ")
	if len(token) < 2 {
		errorResponse(c, errors.UnauthorizedError(""))
	}
	claims, err := m.authAdapter.GetClaim(token[1])
	if err != nil {
		errorResponse(c, err)
		return
	}
	uid := claims[config.UserIdClaim]
	user, err := m.userRepo.FindAuthUser(uid)
	if err != nil {
		errorResponse(c, errors.UnauthorizedError("User not found"))
		return
	}

	userManager := func() interface{} {
		return user
	}

	c.Set(config.AuthUserContextKey, userManager)
	c.Next()
}

func errorResponse(ctx *gin.Context, err errors.ApplicationError) {
	code, message := err.Resolve()
	ctx.JSON(code, gin.H{
		"status": false,
		"error":  message,
	})
	ctx.Abort()
}
