package types

const (
	MEMBER_STATUS_REJECTED  = 0
	MEMBER_STATUS_APPROVED  = 1
	MEMBER_STATUS_INVITED   = 2
	MEMBER_STATUS_REQUESTED = 3
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
		if m.Status == MEMBER_STATUS_APPROVED {
			res = append(res, m)
		}
	}
	return res
}

func (dm DashMembers) ByUserId() map[int]*DashMember {
	res := make(map[int]*DashMember, len(dm))
	for _, v := range dm {
		res[v.UserId] = v
	}
	return res
}

func (dm DashMembers) ByDashId() map[int]*DashMember {
	res := make(map[int]*DashMember, len(dm))
	for _, v := range dm {
		res[v.DashId] = v
	}
	return res
}

func (dm DashMembers) HasUser(userId int) *DashMember {
	for _, v := range dm {
		if v.UserId == userId {
			return v
		}
	}
	return nil
}

func (dm DashMembers) HasDash(dashId int) *DashMember {
	for _, v := range dm {
		if v.DashId == dashId {
			return v
		}
	}
	return nil
}
