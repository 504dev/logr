package controllers

import (
	"github.com/504dev/logr/repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminController struct {
	repos *repo.Repos
}

func NewAdminController(repos *repo.Repos) *AdminController {
	return &AdminController{
		repos: repos,
	}
}

func (adm *AdminController) Users(c *gin.Context) {
	users, _ := adm.repos.User.GetAll()
	c.JSON(http.StatusOK, users)
}

func (adm *AdminController) UserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	usr, _ := adm.repos.User.GetById(id)
	c.JSON(http.StatusOK, usr)
}

func (adm *AdminController) Dashboards(c *gin.Context) {
	dashboards, _ := adm.repos.Dashboard.GetAll()
	c.JSON(http.StatusOK, dashboards)
}

func (adm *AdminController) DashboardById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	dash, _ := adm.repos.Dashboard.GetById(id)
	c.JSON(http.StatusOK, dash)
}
