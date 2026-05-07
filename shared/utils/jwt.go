package utils

import (
	"errors"
	"nutriz-backend-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	clockSkew       = 30 * time.Second
)

type JwtClaims struct {
	IdUser string `json:"id_user"`
	jwt.RegisteredClaims
}

func jwtMethod() jwt.SigningMethod {
	return jwt.SigningMethodHS256
}

func CreateJwtToken(cfg *config.Env, claims JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwtMethod(), claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

func DecodeJwtToken(cfg *config.Env, tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwtMethod().Alg() {
				return nil, ErrInvalidToken
			}
			return []byte(cfg.JWT.Secret), nil
		},
		jwt.WithLeeway(clockSkew),
	)

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid || claims.IdUser == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
