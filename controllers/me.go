package controllers

import (
	"github.com/504dev/kidlog/models/dashboard"
	"github.com/504dev/kidlog/models/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MeController struct{}

func (_ MeController) Me(c *gin.Context) {
	id := c.GetInt("userId")
	usr, _ := user.GetById(id)
	c.JSON(http.StatusOK, usr)
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
	dash, err := dashboard.CreateDashboard(id, body.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dash)
}
