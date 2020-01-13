package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"net/http"
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
var STATE = "secretstring"
var REDIRECT_URL = "http://localhost:8080/"

type OAuthController struct{}

func (u OAuthController) Authorize(c *gin.Context) {
	authorizeUrl := conf.AuthCodeURL(STATE)
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
	c.Abort()
}

func (u OAuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != STATE {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}

	tok, err := conf.Exchange(c, code)

	fmt.Println(state, tok, err)

	client := github.NewClient(conf.Client(c, tok))
	repos, _, err := client.Users.Get(c, "")

	fmt.Println(repos, err)

	c.Redirect(http.StatusMovedPermanently, REDIRECT_URL)
	c.Abort()
}
