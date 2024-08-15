package interfaces

import (
	"database/sql"
	"github.com/504dev/logr/types"
)

type DashboardKeyRepo interface {
	GetById(id int) (*types.DashKey, error)
	GetByPubCached(pub string) (*types.DashKey, error)
	GetByIds(ids []int) (types.DashKeys, error)
	GetByDashId(id int) (types.DashKeys, error)
	Create(key *types.DashKey, tx *sql.Tx) error
}
