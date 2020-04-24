package types

type DashMember struct {
	Id     int   `db:"id"      json:"id"`
	DashId int   `db:"dash_id" json:"dash_id"`
	UserId int   `db:"user_id" json:"user_id"`
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
