package server

import (
	"github.com/504dev/kidlog/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "PUT", "POST"},
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
	}

	logsController := controllers.LogsController{}
	{
		r.GET("/logs", auth.EnsureJWT, logsController.Find)
		r.GET("/logs/stats", auth.EnsureJWT, logsController.Stats)
	}

	adminController := controllers.AdminController{}
	{
		r.GET("/dashboards", auth.EnsureJWT, auth.EnsureAdmin, adminController.Dashboards)
		r.GET("/dashboard/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.DashboardById)
		r.GET("/users", auth.EnsureJWT, auth.EnsureAdmin, adminController.Users)
		r.GET("/user/:id", auth.EnsureJWT, auth.EnsureAdmin, adminController.UserById)
	}

	return r
}
