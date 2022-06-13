package controller

import (
	"net/http"

	"github.com/dino16m/golearn-core/config"
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/types"
	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Authenticate(c Validatable) (userId interface{}, err errors.ApplicationError)
}

type JWTAuthService interface {
	GetTokenPair(claim map[string]interface{}) (refreshToken string, authToken string)
	GetToken(claim map[string]interface{}) string
	GetClaim(tokenStr string) (map[string]interface{}, errors.ApplicationError)
}

type RefreshTokenPayload struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type JWTAuthController struct {
	BaseController
	authenticator Authenticator
	authService   JWTAuthService
}

func NewJWTAuthController(authService JWTAuthService, authenticator Authenticator) JWTAuthController {
	return JWTAuthController{authenticator: authenticator, authService: authService}
}

func (ctrl JWTAuthController) RefreshToken(c *gin.Context) {
	var refresh RefreshTokenPayload
	if err := c.ShouldBind(&refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	claim, err := ctrl.authService.GetClaim(refresh.Token)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	if claim["use"] != types.RefreshTokenKey {
		ctrl.ErrorResponse(c, errors.UnauthorizedError("This is not a refresh token"))
		return
	}
	freshClaim := map[string]interface{}{
		config.UserIdClaim: claim[""],
	}
	token := ctrl.authService.GetToken(freshClaim)
	response := map[string]string{
		"token": token,
	}
	ctrl.OkResponse(c, AppResponse{Data: response})
}

func (ctrl JWTAuthController) GetTokenPair(c *gin.Context) {
	userId, err := ctrl.authenticator.Authenticate(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	claim := map[string]interface{}{
		config.UserIdClaim: userId,
	}
	refreshToken, authToken := ctrl.authService.GetTokenPair(claim)
	response := map[string]string{
		"refreshToken": refreshToken,
		"authToken":    authToken,
	}
	ctrl.OkResponse(c, AppResponse{Data: response})
}

func (ctrl JWTAuthController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", ctrl.GetTokenPair)
	router.POST("/refresh-token", ctrl.RefreshToken)
}
