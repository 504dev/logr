package interfaces

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/types"
)

type CountRepo interface {
	Find(filter types.Filter, agg string) (types.Counts, error)
	StatsByDashboardCached(dashId int) ([]*types.DashStatRow, error)
	StatsByLognameCached(dashId int, logname string) ([]*types.DashStatRow, error)
	Store(count *_types.Count) error
	StopBatcher() error
}
