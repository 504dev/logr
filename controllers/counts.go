package controllers

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/count"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CountsController struct{}

func (_ *CountsController) Find(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	hostname := c.Query("hostname")
	version := c.Query("version")
	agg := c.Query("agg")

	if logname == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname required"})
		return
	}

	filter := types.Filter{
		DashId:   dashId,
		Logname:  logname,
		Hostname: hostname,
		Version:  version,
	}

	duration := Logger.Time("response:/counts", time.Millisecond)
	counts, err := count.Find(filter, agg)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	duration()
	Logger.Inc("count:/counts", 1)
	c.JSON(http.StatusOK, counts.Format())
}

func (_ *CountsController) FindSnippet(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	hostname := c.Query("hostname")
	keyname := c.Query("keyname")
	kind := c.Query("kind")
	timestamp, _ := strconv.ParseInt(c.Query("timestamp"), 10, 64)

	if logname == "" || hostname == "" || keyname == "" || kind == "" || timestamp == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname, hostname, kind, keyname and timestamp required"})
		return
	}

	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
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

	counts, err := count.Find(filter, count.AggMinute)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list := counts.Format()
	for _, v := range list {
		if kind == v.Kind {
			c.JSON(http.StatusOK, v)
			return
		}
	}

	c.JSON(http.StatusOK, nil)
}

func (_ *CountsController) Stats(c *gin.Context) {
	dashId := c.GetInt("dashId")
	duration := Logger.Time("response:/counts/stats", time.Millisecond)
	stats, err := count.GetDashStats(dashId)
	if err != nil {
		Logger.Error(err)
	}
	duration()
	c.JSON(http.StatusOK, stats)
}

func (_ *CountsController) Lognames(c *gin.Context) {
	dashId := c.GetInt("dashId")
	duration := Logger.Time("response:/counts/lognames", time.Millisecond)
	stats, err := count.GetDashLognamesCached(dashId)
	if err != nil {
		Logger.Error(err)
	}
	duration()
	c.JSON(http.StatusOK, stats)
}
