package controller

import (
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/event"
	"github.com/gin-gonic/gin"
)

type Validatable interface {
	ShouldBind(obj interface{}) error
}

type PasswordChangeForm struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
}

type UserService interface {
	CreateUser(ctx Validatable) (interface{}, errors.ApplicationError)
	ChangePassword(user interface{}, dto PasswordChangeForm) errors.ApplicationError
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
	user, err := ctrl.userService.CreateUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	ctrl.dispatcher.Dispatch(event.UserCreated, user)
	ctrl.OkResponse(c, AppResponse{Data: user})
}

func (ctrl AuthController) ChangePassword(c *gin.Context) {
	user, err := ctrl.GetAuthUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}

	var dto PasswordChangeForm
	e := c.ShouldBind(&dto)
	if e != nil {
		ctrl.ErrorResponse(c, errors.ValidationError(e.Error()))
		return
	}

	err = ctrl.userService.ChangePassword(user, dto)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	ctrl.OkResponse(c, AppResponse{})

}

func (ctrl AuthController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ctrl.Signup)
	router.POST("/change-password", ctrl.ChangePassword)
}
