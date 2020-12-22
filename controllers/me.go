package controllers

import (
	"fmt"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/models/dashboard"
	"github.com/504dev/logr/models/dashkey"
	"github.com/504dev/logr/models/dashmember"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
)

type MeController struct{}

func (_ *MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, _ := user.GetById(id)
	c.JSON(http.StatusOK, usr)
}
func (_ *MeController) ShareDashboard(c *gin.Context) {
	ownerId := c.GetInt("userId")
	dashId := c.GetInt("dashId")
	dash, _ := dashboard.GetById(dashId)
	if ownerId != dash.OwnerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	username := c.Query("username")
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
	err := dashmember.Create(&membership)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	membership.User = userTo
	c.JSON(http.StatusOK, membership)
}

func (_ *MeController) RemoveMember(c *gin.Context) {
	userId := c.GetInt("userId")
	dashId := c.GetInt("dashId")
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dash, _ := dashboard.GetById(dashId)
	if userId != dash.OwnerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	err := dashmember.Remove(id)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, id)
}

func (_ *MeController) DashboardsOwn(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := dashboard.GetUserDashboards(userId)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, dash := range dashboards {
		dash.Keys, _ = dashkey.GetByDashId(dash.Id)
		dash.Owner, _ = user.GetById(dash.OwnerId)
		dash.Members, _ = dashmember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = user.GetById(member.UserId)
		}
	}
	c.JSON(http.StatusOK, dashboards)
}

func (_ *MeController) DashboardsShared(c *gin.Context) {
	userId := c.GetInt("userId")
	shared, err := dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, dash := range shared {
		dash.Owner, _ = user.GetById(dash.OwnerId)
		dash.Members, _ = dashmember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = user.GetById(member.UserId)
		}
	}
	c.JSON(http.StatusOK, shared)
}

func (_ *MeController) Dashboards(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := dashboard.GetUserDashboards(userId)
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for _, dash := range dashboards {
		dash.Keys, _ = dashkey.GetByDashId(dash.Id)
	}
	shared, err := dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	dashboards = append(dashboards, shared...)
	for _, dash := range dashboards {
		dash.Owner, _ = user.GetById(dash.OwnerId)
		dash.Members, _ = dashmember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = user.GetById(member.UserId)
		}
	}

	c.JSON(http.StatusOK, dashboards)
}

func (_ *MeController) AddDashboard(c *gin.Context) {
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

	dash.Owner, _ = user.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (_ *MeController) DashRequired(name string) func(c *gin.Context) {
	return func(c *gin.Context) {
		dashId, _ := strconv.Atoi(c.Param(name))
		if dashId == 0 {
			dashId, _ = strconv.Atoi(c.Query(name))
		}
		if dashId == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("%v required", name)})
			return
		}
		dash, err := dashboard.GetById(dashId)
		if err != nil {
			Logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if dash == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Set("dashId", dashId)
		c.Set("dash", dash)
	}
}

func (_ *MeController) MyDash(c *gin.Context) {
	tmp, _ := c.Get("dash")
	dash := tmp.(*types.Dashboard)
	ownerId := c.GetInt("userId")
	if dash.OwnerId != ownerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}

func (_ *MeController) MyDashOrShared(c *gin.Context) {
	tmp, _ := c.Get("dash")
	dash := tmp.(*types.Dashboard)
	userId := c.GetInt("userId")
	role := c.GetInt("role")
	if dash.OwnerId != userId {
		systemIds := dashboard.GetSystemIds(role)
		if sort.SearchInts(systemIds, dash.Id) == len(systemIds) {
			members, err := dashmember.GetAllByUserId(userId)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if members.ApprovedOnly().HasDash(dash.Id) == nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
	}
}

func (_ *MeController) EditDashboard(c *gin.Context) {
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

	dash.Owner, _ = user.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (_ *MeController) DeleteDashboard(c *gin.Context) {
	err := dashboard.Delete(c.GetInt("dashId"))
	if err != nil {
		Logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, c.GetInt("dash"))
}
