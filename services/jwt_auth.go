package services

import (
	"fmt"
	"time"

	"github.com/dino16m/golearn-core/config"
	"github.com/dino16m/golearn-core/errors"
	"github.com/dino16m/golearn-core/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClaims = map[string]interface{}

type RefreshValidator interface {
	BlackListFamily(family string)
	InvalidateJTI(jti string)
	ValidateJTI(family string, jti string) bool
}

type JWTAuthService struct {
	options          config.JwtOptions
	refreshValidator RefreshValidator
}

type TokenPair struct {
	Refresh string `json:"refreshToken"`
	Auth    string `json:"authToken"`
}

func NewJWTAuthService(options config.JwtOptions) JWTAuthService {
	return JWTAuthService{options: options}
}

func (a JWTAuthService) RefreshToken(refreshToken string) (TokenPair, errors.ApplicationError) {
	claim, err := a.GetClaim(refreshToken)

	if err != nil {
		return TokenPair{}, err
	}

	if claim["use"] != types.RefreshTokenKey {
		return TokenPair{}, errors.UnauthorizedError("This is not a refresh token")
	}

	jti := claim["jti"].(string)

	family := claim["fam"].(string)

	valid := a.refreshValidator.ValidateJTI(family, jti)
	if !valid {
		a.refreshValidator.BlackListFamily(family)
		return TokenPair{}, errors.UnauthorizedError("Blacklisted token used")
	}

	a.refreshValidator.InvalidateJTI(jti)

	freshClaim := JWTClaims{
		config.UserIdClaim: claim[config.UserIdClaim],
	}

	refresh := a.GetRefreshToken(jwt.MapClaims{"uid": freshClaim[config.UserIdClaim], "fam": family})
	auth := a.GetToken(claim)

	return TokenPair{
		Refresh: refresh,
		Auth:    auth,
	}, nil
}

func (a JWTAuthService) GetRefreshToken(baseClaims JWTClaims) string {
	baseClaims["iat"] = time.Now().Unix()
	baseClaims["nbf"] = time.Now().Unix()
	baseClaims["exp"] = time.Now().Add(a.options.MaxRefresh).Unix()
	baseClaims["use"] = types.RefreshTokenKey
	baseClaims["jti"] = getJTI()

	if baseClaims["fam"] == nil || baseClaims["fam"] == "" {
		baseClaims["fam"] = getJTI()
	}

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

func getJTI() string {
	id := uuid.New()
	return id.String()
}
