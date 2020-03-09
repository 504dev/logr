package types

type Counter struct {
	DashId    int    `db:"dash_id"   json:"dash_id"`
	Timestamp int64  `db:"timestamp" json:"timestamp"`
	Hostname  string `db:"hostname"  json:"hostname"`
	Name      string `db:"name"      json:"name"`
	*Inc
	*Avg
	*Max
	*Min
	*Per
}

type Avg struct {
	Sum float64 `json:"avg_sum"`
	Num int     `json:"avg_num"`
}

type Inc struct {
	Inc float64 `json:"inc"`
}

type Per struct {
	Taken float64 `json:"per_tkn"`
	Total float64 `json:"per_ttl"`
}

type Max struct {
	Max float64 `json:"max"`
}

type Min struct {
	Min float64 `json:"min"`
}
