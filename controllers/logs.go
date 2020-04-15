package controllers

import (
	"encoding/json"
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/log"
	"github.com/504dev/kidlog/models/ws"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogsController struct{}

func (_ LogsController) Stats(c *gin.Context) {
	dashboards, err := dashboard.GetUserDashboards(c.GetInt("userId"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(dashboards) == 0 {
		c.JSON(http.StatusOK, []int{})
		return
	}
	ids := dashboards.Ids()
	stats, err := log.GetDashStats(ids)
	Logger.Error(err)
	c.JSON(http.StatusOK, stats)
}

func (_ LogsController) Pause(c *gin.Context) {
	userId := c.GetInt("userId")
	sockId := c.Query("sock_id")
	state := false
	if c.Query("state") == "true" {
		state = true
	}
	ws.SockMap.SetPaused(userId, sockId, state)
	c.Status(http.StatusOK)
}

func (_ LogsController) Find(c *gin.Context) {
	dashId, _ := strconv.Atoi(c.Query("dash_id"))
	sockId := c.Query("sock_id")
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
		Timestamp: timestamp,
		DashId:    dashId,
		Logname:   logname,
		Hostname:  hostname,
		Level:     level,
		Version:   version,
		Pid:       pid,
		Message:   message,
		Offset:    offset,
		Limit:     limit,
	}
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
