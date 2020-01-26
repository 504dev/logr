package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogsController struct{}

func (_ LogsController) Stats(c *gin.Context) {
	dashId, _ := strconv.Atoi(c.Query("dash_id"))
	if dashId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "dash_id required"})
		return
	}
	stats, err := log.GetDashStats(dashId)
	fmt.Println(err)
	c.JSON(200, stats)
}

func (_ LogsController) Find(c *gin.Context) {
	dashId, _ := strconv.Atoi(c.Query("dash_id"))
	if dashId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "dash_id required"})
		return
	}

	userId := c.GetInt("userId")
	dash, err := dashboard.GetById(dashId)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if dash.OwnerId != userId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	logname := c.Query("logname")
	hostname := c.Query("hostname")
	message := c.Query("message")
	level := c.Query("level")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)

	from, _ := strconv.ParseInt(c.Query("timestamp.from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("timestamp.to"), 10, 64)

	where := log.Filter{
		Timestamp: [2]int64{from, to},
		DashId:    dashId,
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
}
