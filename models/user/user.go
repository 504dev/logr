package user

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
)

func getByField(fieldname string, val interface{}) (*User, error) {
	var user User
	conn := mysql.Conn()
	sql := fmt.Sprintf("SELECT id, github_id, username FROM users WHERE %v = ?", fieldname)
	row := conn.QueryRowx(sql, val)
	err := row.StructScan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAll() (Users, error) {
	conn := mysql.Conn()
	rows, err := conn.Queryx("SELECT id, github_id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dashboards := make(Users, 0)

	for rows.Next() {
		var user User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, &user)
	}
	return dashboards, nil
}

func GetById(id int) (*User, error) {
	return getByField("id", id)
}
