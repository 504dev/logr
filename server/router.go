package server

import (
	"github.com/504dev/kidlog/controllers"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	logsController := controllers.LogsController{}
	r.GET("/logs", logsController.Find)
	r.GET("/dashboards", func(c *gin.Context) {
		dashboards, _ := dashboard.GetAll()
		c.JSON(200, dashboards)
	})
	r.GET("/dashboard/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		dash, _ := dashboard.GetById(id)
		c.JSON(200, dash)
	})
	r.GET("/users", func(c *gin.Context) {
		users, _ := user.GetAll()
		c.JSON(200, users)
	})
	r.GET("/user/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		usr, _ := user.GetById(id)
		c.JSON(200, usr)
	})

	// me
	me := controllers.MeController{}
	{
		r.GET("/me", me.Me)
		r.GET("/me/dashboards", me.Dashboards)
	}

	// oauth
	oauth := controllers.OAuthController{}
	{
		r.GET("/oauth/signin", oauth.Authorize)
		r.GET("/oauth/callback", oauth.Callback)
	}

	return r
}
