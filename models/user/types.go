package user

type User struct {
	Id       int    `db:"id"`
	GithubId int    `db:"github_id"`
	Username string `db:"username"`
}
type Users []*User
