package types

const (
	MemberStatusRejected  = 0
	MemberStatusApproved  = 1
	MemberStatusInvited   = 2
	MemberStatusRequested = 3
)

type DashMember struct {
	Id     int   `db:"id"      json:"id"`
	DashId int   `db:"dash_id" json:"dash_id"`
	UserId int   `db:"user_id" json:"user_id"`
	Status int   `db:"status"  json:"status"`
	User   *User `json:"user"`
}
type DashMembers []*DashMember

func (dm DashMembers) DashIds() []int {
	ids := make([]int, len(dm))
	for k, v := range dm {
		ids[k] = v.DashId
	}
	return ids
}

func (dm DashMembers) ApprovedOnly() DashMembers {
	res := make(DashMembers, 0, len(dm))
	for _, m := range dm {
		if m.Status == MemberStatusApproved {
			res = append(res, m)
		}
	}
	return res
}

func (ds DashMembers) ByUserId() map[int]*DashMember {
	res := make(map[int]*DashMember, len(ds))
	for _, v := range ds {
		res[v.UserId] = v
	}
	return res
}

func (ds DashMembers) ByDashId() map[int]*DashMember {
	res := make(map[int]*DashMember, len(ds))
	for _, v := range ds {
		res[v.DashId] = v
	}
	return res
}
