package controllers

import (
	. "github.com/504dev/kidlog/logger"
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/dashkey"
	"github.com/504dev/kidlog/models/user"
	"github.com/504dev/kidlog/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MeController struct{}

func (_ MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, _ := user.GetById(id)
	c.JSON(http.StatusOK, usr)
}
func (_ MeController) ShareDashboard(c *gin.Context) {
	ownerId := c.GetInt("userId")
	dashId := c.GetInt("dashId")
	dash, _ := dashboard.GetById(dashId)
	if ownerId != dash.OwnerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	username := c.Param("username")
	userTo, _ := user.GetByUsername(username)
	if userTo == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if ownerId == userTo.Id {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	membership := types.DashMember{
		DashId: dash.Id,
		UserId: userTo.Id,
	}
	err := dashboard.AddMember(&membership)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, membership)
}

func (_ MeController) Dashboards(c *gin.Context) {
	id := c.GetInt("userId")
	dashboards, err := dashboard.GetUserDashboards(id)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, dash := range dashboards {
		dash.Keys, err = dashkey.GetByDashId(dash.Id)
	}
	shared, err := dashboard.GetShared(id)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	dashboards = append(dashboards, shared...)

	c.JSON(http.StatusOK, dashboards)
}

func (_ MeController) Shared(c *gin.Context) {
	id := c.GetInt("userId")
	dashboards, err := dashboard.GetShared(id)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, dashboards)
}

func (_ MeController) AddDashboard(c *gin.Context) {
	var dash *types.Dashboard
	if err := c.BindJSON(&dash); err != nil {
		return
	}

	if dash.Name == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	dash.OwnerId = c.GetInt("userId")
	err := dashboard.Create(dash)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dash)
}

func (_ MeController) IsMyDash(c *gin.Context) {
	ownerId := c.GetInt("userId")
	dashId, _ := strconv.Atoi(c.Param("dashid"))
	if dashId == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dash, err := dashboard.GetById(dashId)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if dash == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if dash.OwnerId != ownerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Set("dashId", dashId)
	c.Set("dash", dash)
}

func (_ MeController) EditDashboard(c *gin.Context) {
	var dash *types.Dashboard
	if err := c.BindJSON(&dash); err != nil {
		return
	}

	dash.Id = c.GetInt("dashId")

	err := dashboard.Update(dash)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dash)
}

func (_ MeController) DeleteDashboard(c *gin.Context) {
	err := dashboard.Delete(c.GetInt("dashId"))
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, c.GetInt("dash"))
}
