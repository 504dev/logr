package types

import "time"

type User struct {
	Id        int        `db:"id"         json:"id"`
	GithubId  int64      `db:"github_id"  json:"github_id"`
	Username  string     `db:"username"   json:"username"`
	Role      int        `db:"role"       json:"role"`
	LoginAt   *time.Time `db:"login_at"   json:"login_at"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}

type Users []*User

const (
	RoleAdmin  = 1
	RoleUser   = 2
	RoleViewer = 3
	RoleDemo   = 4
)

const (
	UserLogrId = 1
	UserDemoId = 2
)
