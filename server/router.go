package server

import (
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
		res := map[string]string{
			"version": Logger.GetVersion(),
			"org":     config.Get().OAuth.Github.Org,
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/free-token", func(c *gin.Context) {
		usr, err := user.GetById(2)
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
					ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(config.Get().OAuth.JwtSecret))
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
		r.GET("/oauth/callback", auth.Callback)
	}

	// me
	me := controllers.MeController{}
	{
		r.GET("/api/me", auth.EnsureJWT, me.Me)
		r.GET("/api/me/dashboards", auth.EnsureJWT, me.Dashboards)
		r.POST("/api/me/dashboard", auth.EnsureJWT, auth.EnsureUser, me.AddDashboard)
		r.POST(
			"/api/me/dashboard/share/:dash_id/to/:username",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.ShareDashboard,
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
		r.GET("/api/logs/stats/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.Stats)
		r.GET("/api/logs/lognames/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, logsController.Lognames)
	}

	countsController := controllers.CountsController{}
	{
		r.GET("/api/counts", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.Find)
		r.GET("/api/counts/snippet", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.FindSnippet)
		r.GET("/api/counts/stats/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.Stats)
		r.GET("/api/counts/lognames/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDashOrShared, countsController.Lognames)

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

	return r
}
