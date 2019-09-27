package user

type User struct {
	GithubId int    `json:"gid"`
	Username string `json:"username"`
}
type Users []*User
