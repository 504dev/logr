package user

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
)

func findOneByField(fieldname string, val interface{}) (*types.User, error) {
	conn := mysql.Conn()
	user := types.User{}
	sql := fmt.Sprintf("SELECT id, github_id, username, role FROM users WHERE %v = ? LIMIT 1", fieldname)
	err := conn.Get(&user, sql, val)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAll() (types.Users, error) {
	conn := mysql.Conn()
	users := types.Users{}
	err := conn.Select(&users, "SELECT id, github_id, username, role FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetById(id int) (*types.User, error) {
	return findOneByField("id", id)
}

func GetByGithubId(id int64) (*types.User, error) {
	return findOneByField("github_id", id)
}

func Create(githubId int64, username string) (*types.User, error) {
	conn := mysql.Conn()

	sql := "INSERT INTO users (github_id, username, role) VALUES (?, ?, ?)"
	values := []interface{}{githubId, username, types.RoleUser}

	_, err := conn.Exec(sql, values...)
	if err != nil {
		return nil, err
	}

	user, err := GetByGithubId(githubId)

	return user, err
}
