package controllers

import (
	"fmt"
	"github.com/504dev/logr/models/dashboard"
	"github.com/504dev/logr/models/dashkey"
	"github.com/504dev/logr/models/dashmember"
	"github.com/504dev/logr/models/user"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"net/http"
	"sort"
	"strconv"
)

type MeController struct{}

func (_ *MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, err := user.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usr)
}
func (_ *MeController) AddMember(c *gin.Context) {
	idash, _ := c.Get("dash")
	dash := idash.(*types.Dashboard)
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "username required"})
		return
	}
	if username == c.GetString("username") {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "share to owner denied"})
		return
	}
	members, err := dashmember.GetAllByDashId(dash.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	limit := 20
	if len(members) > limit {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "members limit"})
		return
	}
	userTo, err := user.GetByUsername(username)
	if userTo == nil && err == nil {
		client := github.NewClient(nil)
		userGithub, _, err := client.Users.Get(c, username)
		if err == nil {
			created, err := user.Create(*userGithub.ID, username, types.RoleUser)
			if err == nil {
				userTo = created
			}
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if userTo == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "user not found"})
		return
	}
	membership := types.DashMember{
		DashId: dash.Id,
		UserId: userTo.Id,
	}
	err = dashmember.Create(&membership)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	membership.User = userTo
	c.JSON(http.StatusOK, membership)
}

func (_ *MeController) RemoveMember(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id required"})
		return
	}
	err := dashmember.Remove(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}

func (_ *MeController) DashboardsOwn(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := dashboard.GetUserDashboards(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for _, dash := range dashboards {
		dash.Keys, _ = dashkey.GetByDashId(dash.Id)
	}
	shared, err := dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"msg": "name required"})
		return
	}

	dash.OwnerId = c.GetInt("userId")
	err := dashboard.Create(dash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	dash.Owner, _ = user.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (_ *MeController) EditDashboard(c *gin.Context) {
	var dash *types.Dashboard
	if err := c.BindJSON(&dash); err != nil {
		return
	}

	dash.Id = c.GetInt("dashId")

	err := dashboard.Update(dash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	dash.Owner, _ = user.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (_ *MeController) DeleteDashboard(c *gin.Context) {
	err := dashboard.Delete(c.GetInt("dashId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, c.GetInt("dash"))
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
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
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
	idash, _ := c.Get("dash")
	dash := idash.(*types.Dashboard)
	ownerId := c.GetInt("userId")
	if dash.OwnerId != ownerId {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}

func (_ *MeController) MyDashOrShared(c *gin.Context) {
	idash, _ := c.Get("dash")
	dash := idash.(*types.Dashboard)
	userId := c.GetInt("userId")
	role := c.GetInt("role")
	if dash.OwnerId == userId {
		return
	}
	systemIds := dashboard.GetSystemIds(role)
	if sort.SearchInts(systemIds, dash.Id) == len(systemIds) {
		members, err := dashmember.GetAllByUserId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		if members.ApprovedOnly().HasDash(dash.Id) == nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}
