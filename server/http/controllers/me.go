package controllers

import (
	"fmt"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/types"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v29/github"
	"net/http"
	"sort"
	"strconv"
)

type MeController struct {
	repos *repo.Repos
}

func NewMeController(repos *repo.Repos) *MeController {
	return &MeController{
		repos: repos,
	}
}

func (me *MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, err := me.repos.User.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, usr)
}
func (me *MeController) AddMember(c *gin.Context) {
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
	members, err := me.repos.DashboardMember.GetAllByDashId(dash.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	limit := 20
	if len(members) > limit {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "members limit"})
		return
	}
	userTo, err := me.repos.User.GetByUsername(username)
	if userTo == nil && err == nil {
		client := github.NewClient(nil)
		userGithub, _, err := client.Users.Get(c, username)
		if err == nil {
			created, err := me.repos.User.Create(*userGithub.ID, username, types.RoleUser)
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
	err = me.repos.DashboardMember.Create(&membership)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	membership.User = userTo
	c.JSON(http.StatusOK, membership)
}

func (me *MeController) RemoveMember(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id required"})
		return
	}
	err := me.repos.DashboardMember.Remove(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, id)
}

func (me *MeController) DashboardsOwn(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := me.repos.Dashboard.GetUserDashboards(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for _, dash := range dashboards {
		dash.Keys, _ = me.repos.DashboardKey.GetByDashId(dash.Id)
		dash.Owner, _ = me.repos.User.GetById(dash.OwnerId)
		dash.Members, _ = me.repos.DashboardMember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = me.repos.User.GetById(member.UserId)
		}
	}
	c.JSON(http.StatusOK, dashboards)
}

func (me *MeController) DashboardsShared(c *gin.Context) {
	userId := c.GetInt("userId")
	shared, err := me.repos.Dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for _, dash := range shared {
		dash.Owner, _ = me.repos.User.GetById(dash.OwnerId)
		dash.Members, _ = me.repos.DashboardMember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = me.repos.User.GetById(member.UserId)
		}
	}
	c.JSON(http.StatusOK, shared)
}

func (me *MeController) Dashboards(c *gin.Context) {
	userId := c.GetInt("userId")
	dashboards, err := me.repos.Dashboard.GetUserDashboards(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	for _, dash := range dashboards {
		dash.Keys, _ = me.repos.DashboardKey.GetByDashId(dash.Id)
	}
	shared, err := me.repos.Dashboard.GetShared(userId, c.GetInt("role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	dashboards = append(dashboards, shared...)
	for _, dash := range dashboards {
		dash.Owner, _ = me.repos.User.GetById(dash.OwnerId)
		dash.Members, _ = me.repos.DashboardMember.GetAllByDashId(dash.Id)
		for _, member := range dash.Members {
			member.User, _ = me.repos.User.GetById(member.UserId)
		}
	}

	c.JSON(http.StatusOK, dashboards)
}

func (me *MeController) AddDashboard(c *gin.Context) {
	var dash *types.Dashboard
	if err := c.BindJSON(&dash); err != nil {
		return
	}

	if dash.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "name required"})
		return
	}

	dash.OwnerId = c.GetInt("userId")
	err := me.repos.Dashboard.Create(dash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	dash.Owner, _ = me.repos.User.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (me *MeController) EditDashboard(c *gin.Context) {
	var dash *types.Dashboard
	if err := c.BindJSON(&dash); err != nil {
		return
	}

	dash.Id = c.GetInt("dashId")

	err := me.repos.Dashboard.Update(dash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	dash.Owner, _ = me.repos.User.GetById(dash.OwnerId)

	c.JSON(http.StatusOK, dash)
}

func (me *MeController) DeleteDashboard(c *gin.Context) {
	err := me.repos.Dashboard.Delete(c.GetInt("dashId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, c.GetInt("dash"))
}

func (me *MeController) DashRequired(name string) func(c *gin.Context) {
	return func(c *gin.Context) {
		dashId, _ := strconv.Atoi(c.Param(name))
		if dashId == 0 {
			dashId, _ = strconv.Atoi(c.Query(name))
		}
		if dashId == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("%v required", name)})
			return
		}
		dash, err := me.repos.Dashboard.GetById(dashId)
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

func (me *MeController) MyDashOrShared(c *gin.Context) {
	idash, _ := c.Get("dash")
	dash := idash.(*types.Dashboard)
	userId := c.GetInt("userId")
	role := c.GetInt("role")
	if dash.OwnerId == userId {
		return
	}
	systemIds := me.repos.Dashboard.GetSystemIds(role)
	if sort.SearchInts(systemIds, dash.Id) == len(systemIds) {
		members, err := me.repos.DashboardMember.GetAllByUserId(userId)
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
