package controllers

import (
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MeController struct{}

func (_ MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, _ := user.GetById(id)
	c.JSON(http.StatusOK, usr)
}
func (_ MeController) ShareDashboard(c *gin.Context) {
	var body struct {
		DashId   int    `json:"dash_id"`
		Username string `json:"username"`
	}
	if err := c.BindJSON(&body); err != nil {
		return
	}
	ownerId := c.GetInt("userId")
	dash, _ := dashboard.GetById(body.DashId)
	if ownerId != dash.OwnerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	member, _ := user.GetByUsername(body.Username)
	if member == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	membership := types.DashMember{
		DashId: dash.Id,
		UserId: member.Id,
	}
	err := dashboard.Share(&membership)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, membership)
}

func (_ MeController) Dashboards(c *gin.Context) {
	id := c.GetInt("userId")
	dashboards, _ := dashboard.GetUserDashboards(id)
	c.JSON(http.StatusOK, dashboards)
}

func (_ MeController) AddDashboard(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&body); err != nil {
		return
	}

	id := c.GetInt("userId")
	dash, err := dashboard.Create(id, body.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dash)
}
