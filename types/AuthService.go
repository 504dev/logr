package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecretFunc func() string
}

func NewAuthService(jwtSecretFunc func() string) *AuthService {
	return &AuthService{
		jwtSecretFunc: jwtSecretFunc,
	}
}

func (iam *AuthService) Secret() string {
	return iam.jwtSecretFunc()
}

func (iam *AuthService) ParseToken(token string) (*Claims, *jwt.Token, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(iam.jwtSecretFunc()), nil
	})
	return claims, tkn, err
}

func (iam *AuthService) SignToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(iam.jwtSecretFunc()))
}
