package controller

import (
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/event"
	"github.com/gin-gonic/gin"
)

type Validatable interface {
	ShouldBind(obj interface{}) error
}

type UserService interface {
	CreateUser(ctx Validatable) (interface{}, errors.ApplicationError)
	ChangePassword(user interface{}, ctx Validatable)
}

type AuthController struct {
	BaseController
	userService UserService
	dispatcher  event.Dispatcher
}

func NewAuthController(userService UserService, dispatcher event.Dispatcher) AuthController {
	return AuthController{userService: userService, dispatcher: dispatcher}
}

func (ctrl AuthController) Signup(c *gin.Context) {
	userDTO, err := ctrl.userService.CreateUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}

	ctrl.OkResponse(c, AppResponse{Data: userDTO})
}

func (ctrl AuthController) ChangePassword(c *gin.Context) {
	user, err := ctrl.GetAuthUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
	}
	ctrl.dispatcher.Dispatch(event.UserCreated, user)
	ctrl.userService.ChangePassword(user, c)
}

func (ctrl AuthController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ctrl.Signup)
	router.POST("/change-password", ctrl.ChangePassword)
}
