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
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuthController struct {
	*oauth2.Config
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
}

func (a *AuthController) Authorize(c *gin.Context) {
	authorizeUrl := a.Config.AuthCodeURL(config.Get().OAuth.StateSecret)
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
	c.Abort()
}

func (a *AuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != config.Get().OAuth.StateSecret {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}

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

	REDIRECT_URL := config.Get().OAuth.RedirectUrl
	u, _ := url.Parse(REDIRECT_URL)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "jwt-token",
		Value:  tokenString,
		Path:   "/",
		Domain: u.Hostname(),
		MaxAge: JWT_LIFETIME,
	})

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%v%v", REDIRECT_URL, tokenString))
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
