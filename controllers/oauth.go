package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

const DEFAULT_EXPIRE_TIME = 8 * time.Hour

type AuthController struct {
	*oauth2.Config
	*types.States
}

func (a *AuthController) Init() {
	credentials := config.Get().OAuth.Github
	scopes := []string{"read:user", "user:email", "read:org"}
	a.Config = &oauth2.Config{
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	a.States = types.States{}.Init()
}

func (a *AuthController) Authorize(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	login := c.Query("login")
	callback := c.Query("callback")
	state := a.States.Push(callback)
	authorizeUrl := a.Config.AuthCodeURL(state)
	if login != "" {
		authorizeUrl += "&login=" + login
	}
	c.Redirect(http.StatusMovedPermanently, authorizeUrl)
}

func (a *AuthController) NeedSetup(c *gin.Context) {
	if !config.Get().NeedSetup() {
		c.AbortWithStatus(http.StatusForbidden)
	}
}
func (a *AuthController) Setup(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	callback := c.Query("callback")
	a.States.Set("setup", callback)
	c.Redirect(http.StatusTemporaryRedirect, "https://github.com/settings/apps/new")
}

func (a *AuthController) SetupCallback(c *gin.Context) {
	code := c.Query("code")
	url := fmt.Sprintf("https://api.github.com/app-manifests/%v/conversions", code)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		c.AbortWithStatus(resp.StatusCode)
		return
	}

	var data struct {
		Name         string `json:"name"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Owner        struct {
			ID    int64  `json:"id"`
			Login string `json:"login"`
		} `json:"owner"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = user.Create(data.Owner.ID, data.Owner.Login, types.RoleAdmin)
	if err != nil {
		Logger.Error(err)
	}
	config.Set(func(conf *config.ConfigData) {
		conf.OAuth.Github.ClientId = data.ClientId
		conf.OAuth.Github.ClientSecret = data.ClientSecret
		conf.OAuth.JwtSecret = conf.GetJwtSecret()
	})
	err = config.Save()
	if err != nil {
		Logger.Error(err)
	}
	defer a.Init()
	callback, ok := a.States.Pop("setup")
	if !ok {
		c.JSON(http.StatusOK, data)
	} else {
		c.Redirect(http.StatusMovedPermanently, callback)
	}
}

func (a *AuthController) AuthorizeCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	callback, ok := a.States.Pop(state)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "incorrect state"})
		return
	}

	githubPermit, err := a.Config.Exchange(c, code)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := github.NewClient(a.Config.Client(c, githubPermit))
	userGithub, _, err := client.Users.Get(c, "")
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	Logger.Info(userGithub)
	if org := config.Get().OAuth.Github.Org; org != "" {
		orgs, _, err := client.Organizations.List(c, *userGithub.Login, nil)
		if err != nil {
			Logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		Logger.Info(orgs)
		access := false
		for _, v := range orgs {
			if *v.Login == org {
				access = true
				break
			}
		}
		if !access {
			msg := fmt.Sprintf("membership in the %s organization is required.", org)
			if callback != "" {
				url := fmt.Sprintf("%verror?msg=%v", callback, msg)
				c.Redirect(http.StatusMovedPermanently, url)
			} else {
				c.JSON(http.StatusForbidden, gin.H{"msg": msg})
			}
			return
		}
	}

	userDb, err := user.Upsert(*userGithub.ID, *userGithub.Login, types.RoleUser)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	Logger.Debug(userDb)
	_ = user.LoginAt(userDb.Id)

	var expiresAt time.Time
	if githubPermit.Expiry.IsZero() {
		expiresAt = time.Now().Add(DEFAULT_EXPIRE_TIME)
	} else {
		expiresAt = githubPermit.Expiry
	}

	claims := types.Claims{
		Id:          userDb.Id,
		Role:        userDb.Role,
		GihubId:     *userGithub.ID,
		Username:    *userGithub.Login,
		AccessToken: githubPermit.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	err = claims.EncryptAccessToken(config.Get().GetJwtSecret())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(config.Get().GetJwtSecret()))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if callback != "" {
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%v%v", callback, tokenSigned))
	} else {
		c.JSON(http.StatusOK, tokenSigned)
	}
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
	tkn, err := claims.ParseWithClaims(token, config.Get().GetJwtSecret())

	if err != nil || !tkn.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//
	if claims.AccessTokenCipher != "" {
		err = claims.DecryptAccessToken(config.Get().GetJwtSecret())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	c.Set("claims", claims)
	c.Set("token", token)
	c.Set("userId", claims.Id)
	c.Set("username", claims.Username)
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

func (_ *AuthController) EnsureUser(c *gin.Context) {
	role := c.GetInt("role")
	if role > types.RoleUser {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
