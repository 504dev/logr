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
	oauth := controllers.OAuthController{}
	{
		r.GET("/oauth/authorize", oauth.Authorize)
		r.GET("/oauth/callback", oauth.Callback)
	}

	// me
	me := controllers.MeController{}
	{
		r.GET("/me", oauth.EnsureJWT, me.Me)
		r.GET("/me/dashboards", oauth.EnsureJWT, me.Dashboards)
		r.POST("/me/dashboard", oauth.EnsureJWT, me.AddDashboard)
	}

	logsController := controllers.LogsController{}
	{
		r.GET("/logs", oauth.EnsureJWT, logsController.Find)
	}

	adminController := controllers.AdminController{}
	{
		r.GET("/dashboards", oauth.EnsureJWT, oauth.EnsureAdmin, adminController.Dashboards)
		r.GET("/dashboard/:id", oauth.EnsureJWT, oauth.EnsureAdmin, adminController.DashboardById)
		r.GET("/users", oauth.EnsureJWT, oauth.EnsureAdmin, adminController.Users)
		r.GET("/user/:id", oauth.EnsureJWT, oauth.EnsureAdmin, adminController.UserById)
	}

	return r
}
