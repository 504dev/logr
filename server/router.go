package server

import (
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/controllers"
	. "github.com/504dev/logr/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
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
		r.POST("/api/me/dashboard", auth.EnsureJWT, me.AddDashboard)
		r.POST("/api/me/dashboard/share/:dash_id/to/:username", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.ShareDashboard)
		r.PUT("/api/me/dashboard/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.EditDashboard)
		r.DELETE("/api/me/dashboard/:dash_id", auth.EnsureJWT, me.DashRequired("dash_id"), me.MyDash, me.DeleteDashboard)
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
