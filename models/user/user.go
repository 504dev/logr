package user

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
)

func getByField(fieldname string, val interface{}) (*types.User, error) {
	var user types.User
	conn := mysql.Conn()
	sql := fmt.Sprintf("SELECT id, github_id, username, role FROM users WHERE %v = ?", fieldname)
	row := conn.QueryRowx(sql, val)
	err := row.StructScan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAll() (types.Users, error) {
	conn := mysql.Conn()
	rows, err := conn.Queryx("SELECT id, github_id, username, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dashboards := make(types.Users, 0)

	for rows.Next() {
		var user types.User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, &user)
	}
	return dashboards, nil
}

func GetById(id int) (*types.User, error) {
	return getByField("id", id)
}

func GetByGithubId(id int64) (*types.User, error) {
	return getByField("github_id", id)
}

func Create(githubId int64, username string) (*types.User, error) {
	values := []interface{}{githubId, username, types.Roles["USER"]}
	conn := mysql.Conn()

	sqlstr := `INSERT INTO users (github_id, username, role) VALUES (?, ?, ?)`
	stmt, err := conn.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return nil, err
	}

	user, err := GetByGithubId(githubId)

	return user, err
}
