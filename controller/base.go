package controller

import (
	"fmt"
	"net/http"

	"github.com/dino16m/golearn-core/config"
	"github.com/dino16m/golearn-core/errors"
	"github.com/gin-gonic/gin"
)

type AuthUserManager interface {
	GetAuthUser(c *gin.Context) (interface{}, errors.AppError)
}

// A BaseController is the
// base struct implementing behaviour universal to controllers
type BaseController struct{}

type AppResponse struct {
	Code int
	Data interface{}
}

type UserManager func() any

type Renderable interface {
	Render() map[string]interface{}
}

func getZero[T any]() T {
	var result T
	return result
}
func GetUser[T any](c *gin.Context) (T, errors.ApplicationError) {

	user, exists := c.Get(config.AuthUserContextKey)
	if !exists {
		return getZero[T](), errors.UnauthorizedError("User not authenticated")
	}

	return user.(T), nil
}

// GetAuthUser returns the authenticated user interface and a nil error
// if such user exists.
// It returns a nil user and an error if the user does not exist or if
// no auth manager was registered.
func (b BaseController) GetAuthUser(c *gin.Context) (interface{}, errors.ApplicationError) {
	return GetUser[any](c)
}

// GetBaseURL return the fully qualified url to the root of the app, from the
// request url
func (b BaseController) GetBaseURL(c *gin.Context) string {
	scheme := "http"
	host := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	return baseURL
}

func (b BaseController) ErrorResponse(ctx *gin.Context, err errors.ApplicationError) {
	code, message := err.Resolve()
	ctx.JSON(code, gin.H{
		"status": false,
		"error":  message,
	})
}

func (b BaseController) OkResponse(ctx *gin.Context, res AppResponse) {
	var code = res.Code
	if res.Code == 0 {
		code = http.StatusOK
	}

	data := res.Data
	renderable, ok := data.(Renderable)
	if ok {
		rendered := renderable.Render()
		rendered["status"] = true
		ctx.JSON(code, rendered)
	} else {
		ctx.JSON(code, gin.H{
			"status": true,
			"data":   res.Data,
		})
	}
}
