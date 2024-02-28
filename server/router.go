package server

import (
	"encoding/json"
	"errors"
	"github.com/504dev/logr-go-client/utils"
	"github.com/504dev/logr/cachify"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/controllers"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowAllOrigins: true,
	}))

	r.GET("/api/globals", func(c *gin.Context) {

		globals := map[string]interface{}{
			"version": Logger.GetVersion(),
			"org":     config.Get().OAuth.Github.Org,
			"setup":   config.Get().NeedSetup(),
		}
		if wd, err := os.Getwd(); err == nil {
			wd = wd + "/frontend"
			version := utils.ReadGitTagDir(wd)
			if version == "" {
				version = utils.ReadGitCommitDir(wd)
				if len(version) >= 6 {
					version = version[0:6]
				}
			}
			globals["frontend"] = version
		}
		c.JSON(http.StatusOK, globals)
	})

	r.GET("/api/free-token", func(c *gin.Context) {
		usr, err := user.GetById(types.UserDemoId)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		tokenString, err := cachify.Cachify("free-token", func() (interface{}, error) {
			claims := types.Claims{
				Id:       usr.Id,
				Role:     usr.Role,
				GihubId:  usr.GithubId,
				Username: usr.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(config.Get().GetJwtSecret()))
			if err != nil {
				return nil, err
			}
			return tokenString, err
		}, 4*time.Minute)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, tokenString)
	})

	// oauth
	auth := controllers.AuthController{}
	auth.Init()
	{
		r.GET("/oauth/authorize", auth.Authorize)
		r.GET("/oauth/authorize/callback", auth.AuthorizeCallback)
		r.POST("/oauth/setup", auth.NeedSetup, auth.Setup)
		r.GET("/oauth/setup/callback", auth.NeedSetup, auth.SetupCallback)
	}

	// me
	me := controllers.MeController{}
	{
		r.GET("/api/me", auth.EnsureJWT, me.Me)
		r.GET("/api/me/dashboards", auth.EnsureJWT, me.Dashboards)
		r.POST("/api/me/dashboard", auth.EnsureJWT, auth.EnsureUser, me.AddDashboard)
		r.POST(
			"/api/me/dashboard/:dash_id/member",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.AddMember,
		)
		r.DELETE(
			"/api/me/dashboard/:dash_id/member",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.RemoveMember,
		)
		r.PUT(
			"/api/me/dashboard/:dash_id",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.EditDashboard,
		)
		r.DELETE(
			"/api/me/dashboard/:dash_id",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.DeleteDashboard,
		)
	}

	logsController := controllers.LogsController{}
	{
		r.GET("/api/logs", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.Find)
		r.GET("/api/logs/:dash_id/lognames", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.StatsByDashboard)
		r.GET("/api/logs/:dash_id/stats", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.StatsByLogname)
	}

	countsController := controllers.CountsController{}
	{
		r.GET("/api/counts", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.Find)
		r.GET("/api/counts/:dash_id/snippet", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.FindSnippet)
		r.GET("/api/counts/:dash_id/lognames", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.StatsByDashboard)
		r.GET("/api/counts/:dash_id/stats", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.StatsByLogname)

	}

	adminController := controllers.AdminController{}
	{
		r.GET("/api/dashboards", auth.EnsureJWT, auth.EnsureAdmin, adminController.Dashboards)
		r.GET("/api/dashboard/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.DashboardById)
		r.GET("/api/users", auth.EnsureJWT, auth.EnsureAdmin, adminController.Users)
		r.GET("/api/user/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.UserById)
	}

	wsController := controllers.WsController{}
	r.GET("/ws", wsController.Index)

	r.POST("/webhook", func(c *gin.Context) {
		requestDump, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			Logger.Error(err)
		}
		Logger.Notice(string(requestDump))
		c.AbortWithStatus(http.StatusOK)
	})

	r.POST("/support", func(c *gin.Context) {
		var data struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Message string `json:"message"`
			Token   string `json:"recaptchaToken,omitempty"`
		}
		err := json.NewDecoder(c.Request.Body).Decode(&data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		verifyData, err := CheckRecaptcha(config.Get().RecaptchaSecret, data.Token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		data.Token = ""
		payload, _ := json.Marshal(data)
		Support.Notice("%v %v", string(payload), verifyData)
		c.AbortWithStatus(http.StatusOK)
	})

	return r
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func CheckRecaptcha(secret, response string) (*SiteVerifyResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return nil, err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return &body, errors.New("unsuccessful recaptcha verify request")
	}

	return &body, nil
}
