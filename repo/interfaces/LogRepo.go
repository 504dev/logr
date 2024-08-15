package interfaces

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
)

type LogRepo interface {
	GetByFilter(filter types.Filter) (_types.Logs, error)
	StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error)
	StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error)
	Store(log *_types.Log) error
	StopBatcher() error
}
