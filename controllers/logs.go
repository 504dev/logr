package controllers

import (
	"encoding/json"
	"github.com/504dev/kidlog/logger"
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
	dashId, _ := strconv.Atoi(c.Query("dash_id"))
	if dashId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "dash_id required"})
		return
	}
	stats, err := log.GetDashStats([]int{dashId})
	logger.Error(err)
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
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 64)

	from, _ := strconv.ParseInt(c.Query("timestamp.from"), 10, 64)
	to, _ := strconv.ParseInt(c.Query("timestamp.to"), 10, 64)

	filter := types.Filter{
		Timestamp: [2]int64{from, to},
		DashId:    dashId,
		Logname:   logname,
		Hostname:  hostname,
		Level:     level,
		Message:   message,
		Offset:    offset,
		Limit:     limit,
	}
	if sockId != "" {
		ws.SockMap.SetFilter(userId, sockId, &filter)
	}
	f, _ := json.Marshal(filter)
	logger.Info(string(f))

	logs, err := log.GetByFilter(filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, logs)
}
