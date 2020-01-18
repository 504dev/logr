package user

type User struct {
	Id       int    `db:"id"        json:"id"`
	GithubId int    `db:"github_id" json:"github_id"`
	Username string `db:"username"  json:"username"`
	Role     int    `db:"role"      json:"role"`
}

type Users []*User

var Roles = map[string]int{
	"ADMIN":  1,
	"USER":   2,
	"VIEWER": 2,
}
