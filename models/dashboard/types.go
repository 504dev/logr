package dashboard

type Dashboard struct {
	Id         int
	Name       string
	PublicKey  string
	PrivateKey string
}
type Dashboards []*Dashboard
