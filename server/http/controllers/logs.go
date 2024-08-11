package controllers

import (
	"encoding/json"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/log"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LogsController struct {
	sockmap *types.SockMap
}

func NewLogsController(sockmap *types.SockMap) *LogsController {
	return &LogsController{
		sockmap: sockmap,
	}
}

func (lc *LogsController) Find(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	userId := ctx.GetInt("userId")

	logname := ctx.Query("logname")
	hostname := ctx.Query("hostname")
	message := ctx.Query("message")
	pattern := ctx.Query("pattern")
	level := ctx.Query("level")
	version := ctx.Query("version")
	pid, _ := strconv.Atoi(ctx.Query("pid"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.ParseInt(ctx.Query("offset"), 10, 64)

	timestamp := [2]int64{0, 0}
	for k, v := range ctx.QueryArray("timestamp[]") {
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
	sockId := ctx.Query("sock_id")
	if sockId != "" {
		lc.sockmap.SetFilter(userId, sockId, &filter)
	}
	f, _ := json.Marshal(filter)
	Logger.Info(string(f))

	duration := Logger.Time("response:/logs", time.Millisecond)
	logs, err := log.GetByFilter(filter)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, logs)
}

func (_ *LogsController) StatsByLogname(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	duration := Logger.Time("response:/logs/stats", time.Millisecond)

	stats, err := log.StatsByLognameCached(dashId, logname)
	if err != nil {
		Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	c.JSON(http.StatusOK, stats)
}

func (_ *LogsController) StatsByDashboard(c *gin.Context) {
	dashId := c.GetInt("dashId")
	duration := Logger.Time("response:/logs/lognames", time.Millisecond)
	stats, err := log.StatsByDashboardCached(dashId)
	if err != nil {
		Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	c.JSON(http.StatusOK, stats)
}
