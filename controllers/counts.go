package controllers

import (
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/count"
	"github.com/504dev/logr/models/dashboard"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CountsController struct{}

func (_ *CountsController) Find(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	hostname := c.Query("hostname")
	agg := c.Query("agg")

	if logname == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname required"})
		return
	}

	counts, err := count.Find(dashId, logname, hostname, agg)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, counts.Format())
}

func (_ *CountsController) Stats(c *gin.Context) {
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
	stats, err := count.GetDashStats(ids)
	Logger.Error(err)
	c.JSON(http.StatusOK, stats)
}
