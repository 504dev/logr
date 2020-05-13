package controllers

import (
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/count"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CountsController struct{}

func (_ *CountsController) Find(c *gin.Context) {
	dashId := c.GetInt("dashId")
	logname := c.Query("logname")
	hostname := c.Query("hostname")

	if logname == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "logname required"})
		return
	}

	logs, err := count.Find(dashId, logname, hostname)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, logs)
}
