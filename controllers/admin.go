package controllers

import (
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminController struct{}

func (_ *AdminController) Users(c *gin.Context) {
	users, _ := user.GetAll()
	c.JSON(http.StatusOK, users)
}

func (_ *AdminController) UserById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	usr, _ := user.GetById(id)
	c.JSON(http.StatusOK, usr)
}

func (_ *AdminController) Dashboards(c *gin.Context) {
	dashboards, _ := dashboard.GetAll()
	c.JSON(http.StatusOK, dashboards)
}

func (_ *AdminController) DashboardById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	dash, _ := dashboard.GetById(id)
	c.JSON(http.StatusOK, dash)
}
