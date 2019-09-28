package user

type User struct {
	Id       int    `json:"id"`
	GithubId int    `json:"gid"`
	Username string `json:"username"`
}
type Users []*User
