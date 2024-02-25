package server

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			Token   string `json:"recaptchaToken"`
		}
		err := json.NewDecoder(c.Request.Body).Decode(&data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		posturl := "https://www.google.com/recaptcha/api/siteverify"
		body := []byte(fmt.Sprintf(`{ "secret": "%s", "response": "%s" }`, config.Get().ReCaptcha, data.Token))

		r, err := http.Post(posturl, "application/json", bytes.NewBuffer(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		var result struct {
			Success bool `json:"success"`
		}
		err = json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if result.Success == false {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		payload, _ := c.GetRawData()
		Support.Info("%v %v", string(payload))
		c.AbortWithStatus(http.StatusOK)
	})

	return r
}
