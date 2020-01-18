package controllers

import (
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/gin-gonic/gin"
)

type MeController struct{}

func (_ MeController) Me(c *gin.Context) {
	id := c.GetInt("id")
	usr, _ := user.GetById(id)
	c.JSON(200, usr)
}
func (_ MeController) Dashboards(c *gin.Context) {
	id := c.GetInt("id")
	dashboards, _ := dashboard.GetUserDashboards(id)
	c.JSON(200, dashboards)
}
