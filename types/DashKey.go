package types

type DashKey struct {
	Id         int    `db:"id"          json:"id"`
	DashId     int    `db:"dash_id"     json:"dash_id"`
	Name       string `db:"name"        json:"name"`
	PublicKey  string `db:"public_key"  json:"public_key"`
	PrivateKey string `db:"private_key" json:"private_key"`
}
type DashKeys []*DashKey
