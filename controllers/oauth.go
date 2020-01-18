package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var conf = &oauth2.Config{
	ClientID:     "a3e0eabef800cd0e7a84",
	ClientSecret: "95344c1682df6e82e71652398dcf9f44b1c6ed8d",
	Scopes:       []string{"user"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

type OAuthController struct{}

func (_ OAuthController) Authorize(c *gin.Context) {
	authorizeUrl := conf.AuthCodeURL(config.Get().OAuth.StateSecret)
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
	c.Abort()
}

func (_ OAuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != config.Get().OAuth.StateSecret {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}

	tok, err := conf.Exchange(c, code)

	client := github.NewClient(conf.Client(c, tok))
	user, _, err := client.Users.Get(c, "")

	// TODO create user if not exist

	fmt.Println(user, err)

	claims := types.Claims{
		GihubId:     *user.ID,
		AccessToken: tok.AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	claims.EncryptAccessToken()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().OAuth.JwtSecret))

	REDIRECT_URL := config.Get().OAuth.RedirectUrl
	u, _ := url.Parse(REDIRECT_URL)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "jwt-token",
		Value:  tokenString,
		Path:   "/",
		Domain: u.Hostname(),
		MaxAge: 60,
	})

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%v%v", REDIRECT_URL, tokenString))
	c.Abort()
}

func (_ OAuthController) EnsureJWT(c *gin.Context) {
	var tknStr string
	splitted := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(splitted) == 2 {
		tknStr = splitted[1]
	}
	if tknStr == "" {
		tknStr = c.Query("token")
	}

	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().OAuth.JwtSecret), nil
	})

	if err != nil || !tkn.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = claims.DecryptAccessToken()

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("jwt", claims)

	c.Next()
}
