package interfaces

import "github.com/504dev/logr/types"

type DashboardMemberRepo interface {
	GetByDashId(id int) (types.DashMembers, error)
	GetByUserId(id int) (types.DashMembers, error)
	Create(m *types.DashMember) error
	Remove(id int) error
}
