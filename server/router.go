package server

import (
	"github.com/504dev/kidlog/controllers"
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/log"
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

	// oauth
	auth := controllers.AuthController{}
	{
		r.GET("/oauth/authorize", auth.Authorize)
		r.GET("/oauth/callback", auth.Callback)
	}

	// me
	me := controllers.MeController{}
	{
		r.GET("/me", auth.EnsureJWT, me.Me)
		r.GET("/me/dashboards", auth.EnsureJWT, me.Dashboards)
		r.POST("/me/dashboard", auth.EnsureJWT, me.AddDashboard)
		r.POST("/me/dashboard/share/:dashid/to/:username", auth.EnsureJWT, me.IsMyDash, me.ShareDashboard)
		r.PUT("/me/dashboard/:dashid", auth.EnsureJWT, me.IsMyDash, me.EditDashboard)
		r.DELETE("/me/dashboard/:dashid", auth.EnsureJWT, me.IsMyDash, me.DeleteDashboard)
	}

	logsController := controllers.LogsController{}
	{
		r.GET("/logs", auth.EnsureJWT, logsController.Find)
		r.GET("/logs/pause", auth.EnsureJWT, logsController.Pause)
		r.GET("/logs/stats", auth.EnsureJWT, logsController.Stats)
		r.GET("/logs/freq", func(c *gin.Context) {
			stats, err := log.GetFrequentDashboards(1000)
			Logger.Error(err)
			c.JSON(http.StatusOK, stats)
		})
	}

	adminController := controllers.AdminController{}
	{
		r.GET("/dashboards", auth.EnsureJWT, auth.EnsureAdmin, adminController.Dashboards)
		r.GET("/dashboard/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.DashboardById)
		r.GET("/users", auth.EnsureJWT, auth.EnsureAdmin, adminController.Users)
		r.GET("/user/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.UserById)
	}

	wsController := controllers.WsController{}
	r.GET("/ws", wsController.Index)

	return r
}
