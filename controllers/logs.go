package controllers

import (
	"fmt"
	"github.com/504dev/kidlog/models/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogsController struct{}

func (u LogsController) Find(c *gin.Context) {
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
}
