package dashboard

type Dashboard struct {
	Id         int
	Name       string
	OwnerId    int
	PublicKey  string
	PrivateKey string
}
type Dashboards []*Dashboard
