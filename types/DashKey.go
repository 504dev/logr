package types

const (
	DASHKEY_SYSTEM_ID = 1
	DASHKEY_DEMO_ID   = 2
)

type DashKey struct {
	Id         int    `db:"id"          json:"id"`
	DashId     int    `db:"dash_id"     json:"dash_id"`
	Name       string `db:"name"        json:"name"`
	PublicKey  string `db:"public_key"  json:"public_key"`
	PrivateKey string `db:"private_key" json:"private_key"`
}
type DashKeys []*DashKey
type DashKeyMap map[int]*DashKey

func (keys DashKeys) Get(id int) *DashKey {
	for _, v := range keys {
		if v.Id == id {
			return v
		}
	}
	return nil
}
