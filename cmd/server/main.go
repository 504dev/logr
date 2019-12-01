package main

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/mysql"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()

	r := gin.Default()
	r.GET("/logs", func(c *gin.Context) {
		where := log.Filter{DashId: 1}
		logs, _ := log.GetAll(where)
		c.JSON(200, logs)
	})
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

	r.Run(config.Get().Bind.Http)
}
