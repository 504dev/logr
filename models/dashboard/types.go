package dashboard

type Dashboard struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	OwnerId    int    `db:"owner_id"`
	PublicKey  string `db:"public_key"`
	PrivateKey string `db:"private_key"`
}
type Dashboards []*Dashboard
