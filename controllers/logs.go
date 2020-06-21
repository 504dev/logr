package controllers

import (
	"encoding/json"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/dashboard"
	"github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogsController struct{}

func (_ *LogsController) Stats(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := dashboard.GetUserDashboards(userId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	shared, err := dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ids := append(dashboards.Ids(), shared.Ids()...)
	if len(ids) == 0 {
		c.JSON(http.StatusOK, []int{})
		return
	}
	stats, err := log.GetDashStats(ids)
	Logger.Error(err)
	c.JSON(http.StatusOK, stats)
}

func (_ *LogsController) Find(c *gin.Context) {
	dashId := c.GetInt("dashId")
	userId := c.GetInt("userId")

	logname := c.Query("logname")
	hostname := c.Query("hostname")
	message := c.Query("message")
	level := c.Query("level")
	version := c.Query("version")
	pid, _ := strconv.Atoi(c.Query("pid"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)

	timestamp := [2]int64{0, 0}
	for k, v := range c.QueryArray("timestamp[]") {
		if k > 1 {
			break
		}
		t, _ := strconv.ParseInt(v, 10, 64)
		timestamp[k] = t
	}

	filter := types.Filter{
		DashId:    dashId,
		Timestamp: timestamp,
		Logname:   logname,
		Hostname:  hostname,
		Level:     level,
		Version:   version,
		Pid:       pid,
		Message:   message,
		Offset:    offset,
		Limit:     limit,
	}
	sockId := c.Query("sock_id")
	if sockId != "" {
		ws.SockMap.SetFilter(userId, sockId, &filter)
	}
	f, _ := json.Marshal(filter)
	Logger.Info(string(f))

	logs, err := log.GetByFilter(filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, logs)
}
