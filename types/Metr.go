package types

type Metr struct {
	DashId    int                `db:"dash_id"   json:"dash_id"`
	Timestamp int64              `db:"timestamp" json:"timestamp"`
	Hostname  string             `db:"hostname"  json:"hostname"`
	Type      string             `db:"type"      json:"type"`
	Name      string             `db:"name"      json:"name"`
	Scores    map[string]float64 `db:"scores"    json:"scores"`
}
