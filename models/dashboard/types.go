package dashboard

type Dashboard struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	OwnerId    int    `db:"owner_id"`
	PublicKey  string `db:"public_key"`
	PrivateKey string `db:"private_key"`
}
type Dashboards []*Dashboard

type DashMember struct {
	Id     int `db:"id"`
	DashId int `db:"dash_id"`
	UserId int `db:"user_id"`
}
type DashMembers []*DashMembers
