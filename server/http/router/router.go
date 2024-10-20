package router

import (
	"github.com/504dev/logr-go-client/utils"
	"github.com/504dev/logr/config"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/repo"
	"github.com/504dev/logr/server/http/controllers"
	"github.com/504dev/logr/types/jwtservice"
	"github.com/504dev/logr/types/sockmap"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func NewRouter(sockMap *sockmap.SockMap, jwtService *jwtservice.JwtService, repos *repo.Repos) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		AllowAllOrigins: true,
	}))

	r.GET("/api/globals", func(c *gin.Context) {
		globals := map[string]interface{}{
			"version": Logger.GetVersion(),
			"org":     config.Get().OAuth.Github.Org,
			"setup":   config.Get().IsSetupRequired(),
		}
		if wd, err := os.Getwd(); err == nil {
			wd += "/frontend"
			version := utils.ReadGitTagDir(wd)
			if version == "" {
				const versionCropLength = 6
				version = utils.ReadGitCommitDir(wd)
				if len(version) >= versionCropLength {
					version = version[0:versionCropLength]
				}
			}
			globals["frontend"] = version
		}
		c.JSON(http.StatusOK, globals)
	})

	demo := controllers.NewDemoController(jwtService, repos)
	{
		r.GET("/api/free-token", demo.FreeToken)
	}

	// oauth
	auth := controllers.NewAuthController(jwtService, repos)
	{
		r.GET("/oauth/authorize", auth.Authorize)
		r.GET("/oauth/authorize/callback", auth.AuthorizeCallback)
		r.POST("/oauth/setup", auth.NeedSetup, auth.Setup)
		r.GET("/oauth/setup/callback", auth.NeedSetup, auth.SetupCallback)
	}

	// me
	me := controllers.NewMeController(repos)
	{
		r.GET("/api/me", auth.EnsureJWT, me.Me)
		r.GET("/api/me/dashboards", auth.EnsureJWT, me.Dashboards)
		r.POST("/api/me/dashboard", auth.EnsureJWT, auth.EnsureUser, me.AddDashboard)
		r.POST(
			"/api/me/dashboard/:dash_id/member",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.AddMember,
		)
		r.DELETE(
			"/api/me/dashboard/:dash_id/member",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.RemoveMember,
		)
		r.PUT(
			"/api/me/dashboard/:dash_id",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.EditDashboard,
		)
		r.DELETE(
			"/api/me/dashboard/:dash_id",
			auth.EnsureJWT,
			auth.EnsureUser,
			me.DashRequired("dash_id"),
			me.MyDash,
			me.DeleteDashboard,
		)
	}

	// logs
	logs := controllers.NewLogsController(sockMap, repos)
	{
		r.GET("/api/logs",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			logs.Find,
		)
		r.GET(
			"/api/logs/:dash_id/lognames",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			logs.StatsByDashboard,
		)
		r.GET(
			"/api/logs/:dash_id/stats",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			logs.StatsByLogname,
		)
	}

	// counts
	counts := controllers.NewCountsController(repos)
	{
		r.GET(
			"/api/counts",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			counts.Find,
		)
		r.GET(
			"/api/counts/:dash_id/snippet",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			counts.FindSnippet,
		)
		r.GET(
			"/api/counts/:dash_id/lognames",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			counts.StatsByDashboard,
		)
		r.GET(
			"/api/counts/:dash_id/stats",
			auth.EnsureJWT,
			me.DashRequired("dash_id"),
			me.MyDashOrShared,
			counts.StatsByLogname,
		)
	}

	// admin
	admin := controllers.AdminController{}
	{
		r.GET("/api/dashboards", auth.EnsureJWT, auth.EnsureAdmin, admin.Dashboards)
		r.GET("/api/dashboard/:id", auth.EnsureJWT, auth.EnsureAdmin, admin.DashboardById)
		r.GET("/api/users", auth.EnsureJWT, auth.EnsureAdmin, admin.Users)
		r.GET("/api/user/:id", auth.EnsureJWT, auth.EnsureAdmin, admin.UserById)
	}

	// GitHub marketplace
	marketplace := controllers.MarketplaceController{}
	{
		r.POST("/webhook", marketplace.Webhook)
		r.POST("/support", marketplace.Support)
	}

	return r
}
