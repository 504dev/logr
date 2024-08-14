package types

const (
	DASHBOARD_SYSTEM_ID = 1
	DASHBOARD_DEMO_ID   = 2
)

type Dashboard struct {
	Id      int         `db:"id"          json:"id"`
	OwnerId int         `db:"owner_id"    json:"owner_id"`
	Name    string      `db:"name"        json:"name"`
	Keys    DashKeys    `json:"keys"`
	Members DashMembers `json:"members"`
	Owner   *User       `json:"owner"`
}
type Dashboards []*Dashboard
type DashMap map[int]*Dashboard

func (ds Dashboards) Ids() []int {
	ids := make([]int, len(ds))
	for k, v := range ds {
		ids[k] = v.Id
	}
	return ids
}
func (ds Dashboards) ByPrimary() DashMap {
	res := make(DashMap, len(ds))
	for _, v := range ds {
		res[v.Id] = v
	}
	return res
}

type DashStatRow struct {
	DashId   int    `db:"dash_id"  json:"dash_id,omitempty"`
	Hostname string `db:"hostname" json:"hostname,omitempty"`
	Logname  string `db:"logname"  json:"logname"`
	Level    string `db:"level"    json:"level,omitempty"`
	Version  string `db:"version"  json:"version,omitempty"`
	Cnt      int    `db:"cnt"      json:"cnt"`
	Updated  int64  `db:"updated"  json:"updated,omitempty"`
	Message  string `db:"message"  json:"message,omitempty"`
}

type DashStatRows []*DashStatRow
