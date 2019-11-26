package dashboard

type Dashboard struct {
	Id         int         `db:"id"          json:"id"`
	OwnerId    int         `db:"owner_id"    json:"owner_id"`
	Name       string      `db:"name"        json:"name"`
	PublicKey  string      `db:"public_key"  json:"public_key"`
	PrivateKey string      `db:"private_key" json:"private_key"`
	Members    DashMembers `json:"members"`
}
type Dashboards []*Dashboard

type DashMember struct {
	Id     int `db:"id"        json:"id"`
	DashId int `db:"dash_id"   json:"dash_id"`
	UserId int `db:"user_id"   json:"user_id"`
}
type DashMembers []*DashMembers
