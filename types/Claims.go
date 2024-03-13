package types

import (
	"encoding/base64"
	"github.com/504dev/logr-go-client/cipher"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id                int    `json:"id"`
	Role              int    `json:"role"`
	GihubId           int64  `json:"github_id"`
	Username          string `json:"username"`
	AccessToken       string `json:"access_token,omitempty"`
	AccessTokenCipher string `json:"access_token_cipher"`
	jwt.StandardClaims
}

func (claims *Claims) ParseWithClaims(token string, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func (claims *Claims) EncryptAccessToken(secret string) error {
	cipherAccessToken, err := cipher.EncryptAes([]byte(claims.AccessToken), []byte(secret))
	if err != nil {
		return err
	}
	claims.AccessTokenCipher = base64.StdEncoding.EncodeToString(cipherAccessToken)
	//fmt.Println("EncryptAccessToken", p.AccessToken, p.AccessTokenCipher)
	claims.AccessToken = ""

	return nil
}

func (claims *Claims) DecryptAccessToken(secret string) error {
	cipherBytes, _ := base64.StdEncoding.DecodeString(claims.AccessTokenCipher)
	accessToken, err := cipher.DecryptAes(cipherBytes, []byte(secret))
	if err != nil {
		return err
	}
	claims.AccessToken = string(accessToken)
	//fmt.Println("DecryptAccessToken", p.AccessTokenCipher, p.AccessToken)
	claims.AccessTokenCipher = ""

	return nil
}
