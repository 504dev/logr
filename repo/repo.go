package repo

import (
	"github.com/504dev/logr/repo/dashboard"
	"github.com/504dev/logr/repo/dashkey"
	"github.com/504dev/logr/repo/dashmember"
	"github.com/504dev/logr/repo/user"
	"sync"
)

type Repos struct {
	User            *user.UserRepo
	Dashboard       *dashboard.DashboardRepo
	DashboardMember *dashmember.DashboardMemberRepo
	DashboardKey    *dashkey.DashboardKeyRepo
}

var once sync.Once
var repos *Repos

func GetRepos() *Repos {
	once.Do(func() {
		dashboardMemberRepo := dashmember.NewDashboardMemberRepo()
		dashboardKeyRepo := dashkey.NewDashboardKeyRepo()
		repos = &Repos{
			User:            user.NewUserRepo(),
			Dashboard:       dashboard.NewDashboardRepo(dashboardMemberRepo, dashboardKeyRepo),
			DashboardMember: dashboardMemberRepo,
			DashboardKey:    dashboardKeyRepo,
		}
	})
	return repos
}
