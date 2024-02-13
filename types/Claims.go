package types

import (
	"encoding/base64"
	"github.com/504dev/logr-go-client/cipher"
	"github.com/504dev/logr/config"
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

func (p *Claims) EncryptAccessToken() error {
	cipherAccessToken, err := cipher.EncryptAes([]byte(p.AccessToken), []byte(config.Get().GetJwtSecret()))
	if err != nil {
		return err
	}
	p.AccessTokenCipher = base64.StdEncoding.EncodeToString(cipherAccessToken)
	//fmt.Println("EncryptAccessToken", p.AccessToken, p.AccessTokenCipher)
	p.AccessToken = ""

	return nil
}

func (p *Claims) DecryptAccessToken() error {
	cipherBytes, _ := base64.StdEncoding.DecodeString(p.AccessTokenCipher)
	accessToken, err := cipher.DecryptAes(cipherBytes, []byte(config.Get().GetJwtSecret()))
	if err != nil {
		return err
	}
	p.AccessToken = string(accessToken)
	//fmt.Println("DecryptAccessToken", p.AccessTokenCipher, p.AccessToken)
	p.AccessTokenCipher = ""

	return nil
}
