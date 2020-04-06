package user

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
)

func findAllByField(fieldname string, val interface{}, limit int) (types.Users, error) {
	conn := mysql.Conn()
	users := types.Users{}
	sql := fmt.Sprintf("SELECT id, github_id, username, role FROM users WHERE %v = ?", fieldname)
	if limit > 0 {
		sql = fmt.Sprintf("%v LIMIT %v", sql, limit)
	}
	err := conn.Select(&users, sql, val)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func findOneByField(fieldname string, val interface{}) (*types.User, error) {
	users, err := findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
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
