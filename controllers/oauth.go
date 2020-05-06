package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	*oauth2.Config
	states map[string]string
}

func (a *AuthController) Init() {
	conf := config.Get().OAuth.Github
	a.Config = &oauth2.Config{
		ClientID:     conf.ClientId,
		ClientSecret: conf.ClientSecret,
		Scopes:       []string{"user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	a.states = make(map[string]string)
}

func (a *AuthController) Authorize(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	state := fmt.Sprintf("%v_%v", time.Now().Nanosecond(), rand.Int())
	callback := c.Query("callback")
	a.states[state] = callback
	authorizeUrl := a.Config.AuthCodeURL(state)
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
	c.Abort()
}

func (a *AuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	callback, ok := a.states[state]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}
	delete(a.states, state)

	tok, err := a.Config.Exchange(c, code)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := github.NewClient(a.Config.Client(c, tok))
	userGithub, _, err := client.Users.Get(c, "")
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	Logger.Info(userGithub)

	userDb, err := user.GetByGithubId(*userGithub.ID)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if userDb == nil {
		Logger.Error(err)
		userDb, err = user.Create(*userGithub.ID, *userGithub.Login)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	Logger.Debug(userDb)

	JWT_LIFETIME := 60 * 60
	claims := types.Claims{
		Id:          userDb.Id,
		Role:        userDb.Role,
		GihubId:     *userGithub.ID,
		Username:    *userGithub.Login,
		AccessToken: tok.AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(JWT_LIFETIME) * time.Second).Unix(),
		},
	}
	err = claims.EncryptAccessToken()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Get().OAuth.JwtSecret))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if callback != "" {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%v%v", callback, tokenString))
	} else {
		c.JSON(http.StatusOK, tokenString)
	}
	c.Abort()
}

func (_ *AuthController) EnsureJWT(c *gin.Context) {
	var token string
	splitted := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(splitted) == 2 {
		token = splitted[1]
	}
	if token == "" {
		token = c.Query("token")
	}
	if token == "" {
		token = c.PostForm("token")

	}

	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
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

	c.Set("claims", claims)
	c.Set("token", token)
	c.Set("userId", claims.Id)
	c.Set("role", claims.Role)

	c.Next()
}

func (_ *AuthController) EnsureAdmin(c *gin.Context) {
	role := c.GetInt("role")
	if role != types.RoleAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
