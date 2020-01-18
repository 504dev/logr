package types

import (
	"encoding/base64"
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	GihubId           int64  `json:"github_id"`
	AccessToken       string `json:"access_token"`
	AccessTokenCipher string `json:"access_token_cipher"`
	jwt.StandardClaims
}

func (p *Claims) EncryptAccessToken() {
	cipherAccessToken, _ := cipher.EncryptAes([]byte(p.AccessToken), []byte(config.Get().OAuth.JwtSecret))
	p.AccessTokenCipher = base64.StdEncoding.EncodeToString(cipherAccessToken)
	fmt.Println("EncryptAccessToken", p.AccessToken, p.AccessTokenCipher)
	p.AccessToken = ""
}

func (p *Claims) DecryptAccessToken() error {
	cipherBytes, _ := base64.StdEncoding.DecodeString(p.AccessTokenCipher)
	accessToken, err := cipher.DecryptAes(cipherBytes, []byte(config.Get().OAuth.JwtSecret))
	if err != nil {
		return err
	}
	p.AccessToken = string(accessToken)
	fmt.Println("DecryptAccessToken", p.AccessTokenCipher, p.AccessToken)
	p.AccessTokenCipher = ""

	return nil
}
