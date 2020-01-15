package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
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
var STATE = "state-secret"
var REDIRECT_URL = "http://kidlog.loc:8080/jwt/"

type OAuthController struct{}

func (_ OAuthController) Authorize(c *gin.Context) {
	authorizeUrl := conf.AuthCodeURL(STATE)
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
	c.Abort()
}

func (_ OAuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != STATE {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}

	tok, err := conf.Exchange(c, code)

	fmt.Println(state, tok, err)
	fmt.Println(tok.RefreshToken)

	client := github.NewClient(conf.Client(c, tok))
	user, _, err := client.Users.Get(c, "")

	fmt.Println(user, err)
	fmt.Println(c)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.ID,
		"access_token": tok.AccessToken,
	})
	tokenString, err := token.SignedString([]byte("jwt-secret"))

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
