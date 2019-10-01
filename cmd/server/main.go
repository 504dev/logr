package main

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/mysql"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()

	r := gin.Default()
	r.GET("/logs", func(c *gin.Context) {
		c.JSON(200, log.GetAll())
	})
	r.GET("/dashboards", func(c *gin.Context) {
		c.JSON(200, dashboard.GetAll())
	})
	r.GET("/users", func(c *gin.Context) {
		c.JSON(200, user.GetAll())
	})

	r.Run(config.Get().Bind.Http)
}
