package main

import (
	"fmt"
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/mysql"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

type LogHandler struct{}

func (t LogHandler) Write(b []byte) (int, error) {
	return len(b), nil
}

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()

	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, LogHandler{})

	r := gin.Default()
	r.GET("/logs", func(c *gin.Context) {
		dashid, _ := strconv.Atoi(c.Query("dash_id"))
		if dashid == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "day required"})
			return
		}

		logname := c.Query("logname")
		hostname := c.Query("hostname")
		message := c.Query("message")
		level, _ := strconv.Atoi(c.Query("level"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)

		from, _ := strconv.ParseInt(c.Query("timestamp.from"), 10, 64)
		to, _ := strconv.ParseInt(c.Query("timestamp.to"), 10, 64)

		where := log.Filter{
			Timestamp: [2]int64{from, to},
			DashId:    dashid,
			Logname:   logname,
			Hostname:  hostname,
			Level:     level,
			Message:   message,
			Offset:    offset,
			Limit:     limit,
		}
		fmt.Println(where)

		logs, err := log.GetAll(where)
		fmt.Println(err)
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
