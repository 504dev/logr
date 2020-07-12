package types

const (
	DashboardSystemId = 1
	DashboardDemoId   = 2
)

type Dashboard struct {
	Id       int         `db:"id"          json:"id"`
	OwnerId  int         `db:"owner_id"    json:"owner_id"`
	Name     string      `db:"name"        json:"name"`
	Keys     DashKeys    `json:"keys"`
	Members  DashMembers `json:"members"`
	Owner    *User       `json:"owner"`
	Lognames interface{} `json:"lognames,omitempty"`
}
type Dashboards []*Dashboard

func (ds Dashboards) Ids() []int {
	ids := make([]int, len(ds))
	for k, v := range ds {
		ids[k] = v.Id
	}
	return ids
}
func (ds Dashboards) ByPrimary() map[int]*Dashboard {
	res := make(map[int]*Dashboard, len(ds))
	for _, v := range ds {
		res[v.Id] = v
	}
	return res
}

type DashStatRow struct {
	DashId   int    `db:"dash_id"  json:"dash_id"`
	Hostname string `db:"hostname" json:"hostname"`
	Logname  string `db:"logname"  json:"logname"`
	Level    string `db:"level"    json:"level,omitempty"`
	Version  string `db:"version"  json:"version"`
	Cnt      int    `db:"cnt"      json:"cnt"`
	Updated  uint64 `db:"updated"  json:"updated"`
}

type DashStatRows []*DashStatRow
