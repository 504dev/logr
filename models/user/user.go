package user

import (
	"fmt"
	"github.com/504dev/logr/mysql"
	"github.com/504dev/logr/types"
)

func findAllByField(fieldname string, val interface{}, limit int) (types.Users, error) {
	conn := mysql.Conn()
	users := types.Users{}
	sql := fmt.Sprintf("SELECT id, github_id, username, role, login_at, created_at FROM users WHERE %v = ?", fieldname)
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
	err := conn.Select(&users, "SELECT id, github_id, username, role, login_at, created_at FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetById(id int) (*types.User, error) {
	return findOneByField("id", id)
}

func GetByUsername(username string) (*types.User, error) {
	return findOneByField("username", username)
}

func GetByGithubId(id int64) (*types.User, error) {
	return findOneByField("github_id", id)
}

func Upsert(githubId int64, username string, role int) (*types.User, error) {
	user, err := GetByGithubId(githubId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = Create(githubId, username, role)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func Create(githubId int64, username string, role int) (*types.User, error) {
	conn := mysql.Conn()

	sql := `
		INSERT INTO users (github_id, username, role)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE role=VALUES(role)
	`
	values := []interface{}{githubId, username, role}

	_, err := conn.Exec(sql, values...)
	if err != nil {
		return nil, err
	}

	user, err := GetByGithubId(githubId)

	return user, err
}

func LoginAt(id int) error {
	conn := mysql.Conn()
	_, err := conn.Exec("UPDATE users SET login_at = NOW() WHERE id = ?", id)

	return err
}
