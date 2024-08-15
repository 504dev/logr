package interfaces

import "github.com/504dev/logr/types"

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
