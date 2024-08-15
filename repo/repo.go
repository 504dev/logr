package repo

import (
	"github.com/504dev/logr/repo/count"
	"github.com/504dev/logr/repo/dashboard"
	"github.com/504dev/logr/repo/dashkey"
	"github.com/504dev/logr/repo/dashmember"
	"github.com/504dev/logr/repo/interfaces"
	"github.com/504dev/logr/repo/log"
	"github.com/504dev/logr/repo/user"
	"golang.org/x/sync/errgroup"
	"sync"
)

type Repos struct {
	User            interfaces.UserRepo
	Dashboard       interfaces.DashboardRepo
	DashboardMember interfaces.DashboardMemberRepo
	DashboardKey    interfaces.DashboardKeyRepo
	Count           interfaces.CountRepo
	Log             interfaces.LogRepo
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
			Count:           count.NewCountRepo().RunBatcher(),
			Log:             log.NewLogRepo().RunBatcher(),
		}
	})
	return repos
}

func (r *Repos) Stop() error {
	var wg errgroup.Group
	wg.Go(r.Count.StopBatcher)
	wg.Go(r.Log.StopBatcher)
	return wg.Wait()
}
