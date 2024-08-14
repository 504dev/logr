package repo

import (
	"database/sql"
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/repo/count"
	"github.com/504dev/logr/repo/dashboard"
	"github.com/504dev/logr/repo/dashkey"
	"github.com/504dev/logr/repo/dashmember"
	"github.com/504dev/logr/repo/log"
	"github.com/504dev/logr/repo/user"
	"github.com/504dev/logr/types"
	"sync"
)

type UserRepo interface {
	GetAll() (types.Users, error)
	GetById(id int) (*types.User, error)
	GetByUsername(username string) (*types.User, error)
	Upsert(githubId int64, username string, role int) (*types.User, error)
	Create(githubId int64, username string, role int) (*types.User, error)
	LoginAt(id int) error
}

type DashboardRepo interface {
	GetAll() (types.Dashboards, error)
	GetById(id int) (*types.Dashboard, error)
	GetByOwnerId(id int) (types.Dashboards, error)
	GetShared(userId int, role int) (types.Dashboards, error)
	GetSystemIds(role int) []int
	Create(dash *types.Dashboard) error
	Update(dash *types.Dashboard) error
	Remove(id int) error
}

type DashboardMemberRepo interface {
	GetByDashId(id int) (types.DashMembers, error)
	GetByUserId(id int) (types.DashMembers, error)
	Create(m *types.DashMember) error
	Remove(id int) error
}

type DashboardKeyRepo interface {
	GetById(id int) (*types.DashKey, error)
	GetByIds(ids []int) (types.DashKeys, error)
	GetByDashId(id int) (types.DashKeys, error)
	GetByPubCached(pub string) (*types.DashKey, error)
	Create(key *types.DashKey, tx *sql.Tx) error
}

type CountRepo interface {
	Find(filter types.Filter, agg string) (types.Counts, error)
	StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error)
	StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error)
	Store(count *_types.Count) error
	StopQueue() error
}

type LogRepo interface {
	GetByFilter(filter types.Filter) (_types.Logs, error)
	StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error)
	StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error)
	Store(log *_types.Log) error
	StopQueue() error
}

type Repos struct {
	User            UserRepo
	Dashboard       DashboardRepo
	DashboardMember DashboardMemberRepo
	DashboardKey    DashboardKeyRepo
	Count           CountRepo
	Log             LogRepo
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
			Count:           count.NewCountRepo().RunQueue(),
			Log:             log.NewLogRepo().RunQueue(),
		}
	})
	return repos
}
