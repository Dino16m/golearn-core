package services

import (
	"fmt"
	"time"

	"github.com/dino16m/golearn-core/config"
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/types"
	"github.com/golang-jwt/jwt"
)

type JWTClaims = map[string]interface{}

type JWTAuthService struct {
	options config.JwtOptions
}

func NewJWTAuthService(options config.JwtOptions) JWTAuthService {
	return JWTAuthService{options: options}
}

func (a JWTAuthService) GetRefreshToken(baseClaims JWTClaims) string {
	baseClaims["iat"] = time.Now().Unix()
	baseClaims["nbf"] = time.Now().Unix()
	baseClaims["exp"] = time.Now().Add(a.options.MaxRefresh).Unix()
	baseClaims["use"] = types.RefreshTokenKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(baseClaims))

	tokenString, _ := token.SignedString([]byte(a.options.Key))
	return tokenString
}

func (a JWTAuthService) GetToken(claim JWTClaims) string {
	claim["exp"] = time.Now().Add(a.options.Timeout).Unix()
	claim["iat"] = time.Now().Unix()
	claim["nbf"] = time.Now().Unix()
	claim["use"] = types.AuthTokenKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claim))

	tokenString, _ := token.SignedString([]byte(a.options.Key))
	return tokenString
}

func (a JWTAuthService) GetTokenPair(claim JWTClaims) (refreshToken string, authToken string) {
	refreshToken = a.GetRefreshToken(jwt.MapClaims{"uid": claim[config.UserIdClaim]})
	authToken = a.GetToken(claim)
	return refreshToken, authToken
}

func (a JWTAuthService) GetClaim(tokenStr string) (map[string]interface{}, errors.ApplicationError) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalServerError(fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return []byte(a.options.Key), nil
	})
	if err != nil {
		appError, ok := err.(errors.ApplicationError)
		if ok {
			return nil, appError
		} else {
			return nil, errors.InternalServerError(err.Error())
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.UnauthorizedError("Invalid token")
	}

}
