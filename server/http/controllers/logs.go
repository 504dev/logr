package controllers

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types"
	"github.com/504dev/logr/types/sockmap"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LogsController struct {
	sockMap *sockmap.SockMap
	repos   *repo.Repos
}

func NewLogsController(sockMap *sockmap.SockMap, repos *repo.Repos) *LogsController {
	return &LogsController{
		sockMap: sockMap,
		repos:   repos,
	}
}

func (c *LogsController) Find(ctx *gin.Context) {
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

	filter := types.Filter{
		DashId:   dashId,
		Pattern:  pattern,
		Logname:  logname,
		Hostname: hostname,
		Level:    level,
		Version:  version,
		Pid:      pid,
		Message:  message,
		Offset:   offset,
		Limit:    limit,
	}

	if sockId := ctx.Query("sock_id"); sockId != "" {
		c.sockMap.SetFilter(userId, sockId, &filter)
	}
	Logger.Info("Filter %v", filter.String())

	duration := Logger.Time("response:/logs", time.Millisecond)
	logs, err := c.repos.Log.GetByFilter(filter)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, logs)
}

func (c *LogsController) StatsByLogname(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	logname := ctx.Query("logname")
	duration := Logger.Time("response:/logs/stats", time.Millisecond)

	stats, err := c.repos.Log.StatsByLognameCached(dashId, logname)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, stats)
}

func (c *LogsController) StatsByDashboard(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	duration := Logger.Time("response:/logs/lognames", time.Millisecond)
	stats, err := c.repos.Log.StatsByDashboardCached(dashId)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, stats)
}
