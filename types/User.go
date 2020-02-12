package types

type User struct {
	Id       int    `db:"id"        json:"id"`
	GithubId int    `db:"github_id" json:"github_id"`
	Username string `db:"username"  json:"username"`
	Role     int    `db:"role"      json:"role"`
}

type Users []*User

const (
	RoleAdmin  = 1
	RoleUser   = 2
	RoleViewer = 3
)
