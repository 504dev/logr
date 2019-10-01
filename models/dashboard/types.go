package dashboard

type Dashboard struct {
	Id         int         `db:"id"          json:"id"`
	Name       string      `db:"name"        json:"name"`
	OwnerId    int         `db:"owner_id"    json:"owner_id"`
	PublicKey  string      `db:"public_key"  json:"public_key"`
	PrivateKey string      `db:"private_key" json:"private_key"`
	Members    DashMembers `json:"members"`
}
type Dashboards []*Dashboard

type DashMember struct {
	Id       int `db:"id"        json:"id"`
	DashId   int `db:"dash_id"   json:"dash_id"`
	GithubId int `db:"github_id" json:"github_id"`
}
type DashMembers []*DashMembers
