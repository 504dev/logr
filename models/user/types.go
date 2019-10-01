package user

type User struct {
	Id       int    `db:"id"        json:"id"`
	GithubId int    `db:"github_id" json:"github_id"`
	Username string `db:"username"  json:"username"`
}
type Users []*User
