package controllers

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/repo/count"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CountsController struct {
	repos *repo.Repos
}

func NewCountsController(repos *repo.Repos) *CountsController {
	return &CountsController{
		repos: repos,
	}
}

func (c *CountsController) Find(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	logname := ctx.Query("logname")
	hostname := ctx.Query("hostname")
	version := ctx.Query("version")
	agg := ctx.Query("agg")

	if logname == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname required"})
		return
	}

	filter := types.Filter{
		DashId:   dashId,
		Logname:  logname,
		Hostname: hostname,
		Version:  version,
	}

	duration := Logger.Time("response:/counts", time.Millisecond)
	counts, err := c.repos.Count.Find(filter, agg)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, counts.Format())
}

func (c *CountsController) FindSnippet(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	logname := ctx.Query("logname")
	hostname := ctx.Query("hostname")
	keyname := ctx.Query("keyname")
	kind := ctx.Query("kind")
	timestamp, _ := strconv.ParseInt(ctx.Query("timestamp"), 10, 64)

	if logname == "" || hostname == "" || keyname == "" || kind == "" || timestamp == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname, hostname, kind, keyname and timestamp required"})
		return
	}

	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	if limit > 60 {
		limit = 60
	}
	if limit == 0 {
		limit = 15
	}
	from := timestamp - limit*60

	filter := types.Filter{
		DashId:    dashId,
		Timestamp: [2]int64{from, timestamp},
		Logname:   logname,
		Hostname:  hostname,
		Keyname:   keyname,
	}

	counts, err := c.repos.Count.Find(filter, count.AggMinute)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	list := counts.Format()
	for _, v := range list {
		if kind == v.Kind {
			ctx.JSON(http.StatusOK, v)
			return
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *CountsController) StatsByLogname(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	logname := ctx.Query("logname")
	duration := Logger.Time("response:/counts/stats", time.Millisecond)
	stats, err := c.repos.Count.StatsByLognameCached(dashId, logname)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, stats)
}

func (c *CountsController) StatsByDashboard(ctx *gin.Context) {
	dashId := ctx.GetInt("dashId")
	duration := Logger.Time("response:/counts/lognames", time.Millisecond)
	stats, err := c.repos.Count.StatsByDashboardCached(dashId)
	if err != nil {
		Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	duration()
	ctx.JSON(http.StatusOK, stats)
}
