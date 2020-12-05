package controllers

import (
	"encoding/json"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/log"
	"github.com/504dev/logr/models/ws"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LogsController struct{}

func (_ *LogsController) Find(c *gin.Context) {
	dashId := c.GetInt("dashId")
	userId := c.GetInt("userId")

	logname := c.Query("logname")
	hostname := c.Query("hostname")
	message := c.Query("message")
	pattern := c.Query("pattern")
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
		Pattern:   pattern,
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

	duration := Logger.Time("response:/logs", time.Millisecond)
	logs, err := log.GetByFilter(filter)
	if err != nil {
		Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	Logger.Inc("count:/logs", 1)
	c.JSON(http.StatusOK, logs)
}

func (_ *LogsController) Stats(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	duration := Logger.Time("response:/logs/stats", time.Millisecond)

	stats, err := log.GetDashStatsCached(dashId, logname)
	if err != nil {
		Logger.Error(err)
	}
	duration()
	c.JSON(http.StatusOK, stats)
}

func (_ *LogsController) Lognames(c *gin.Context) {
	dashId := c.GetInt("dashId")
	duration := Logger.Time("response:/logs/lognames", time.Millisecond)
	stats, err := log.GetDashLognamesCached(dashId)
	if err != nil {
		Logger.Error(err)
	}
	duration()
	c.JSON(http.StatusOK, stats)
}
