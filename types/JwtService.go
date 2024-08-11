package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secretFunc func() string
}

func NewJwtService(jwtSecretFunc func() string) *JwtService {
	return &JwtService{
		secretFunc: jwtSecretFunc,
	}
}

func (js *JwtService) ParseToken(token string) (*Claims, *jwt.Token, error) {
	secret := js.secretFunc()
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err == nil && claims.AccessTokenCipher != "" {
		err = claims.DecryptAccessToken(secret)
	}
	return claims, tkn, err
}

func (js *JwtService) SignToken(claims *Claims) (string, error) {
	secret := js.secretFunc()
	if err := claims.EncryptAccessToken(secret); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
