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
	ROLE_ADMIN  = 1
	ROLE_USER   = 2
	ROLE_VIEWER = 3
	ROLE_DEMO   = 4
)

const (
	USER_LOGR_ID = 1
	USER_DEMO_ID = 2
)
