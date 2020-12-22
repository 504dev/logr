package dashmember

import (
	"fmt"
	"github.com/504dev/logr/mysql"
	"github.com/504dev/logr/types"
)

func findAllByField(fieldname string, val interface{}, limit int) (types.DashMembers, error) {
	conn := mysql.Conn()
	members := types.DashMembers{}
	sqltext := fmt.Sprintf("SELECT id, dash_id, user_id, status FROM dashboard_members WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	err := conn.Select(&members, sqltext, val)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func findOneByField(fieldname string, val interface{}) (*types.DashMember, error) {
	members, err := findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, nil
	}
	return members[0], nil
}

func GetById(id int) (*types.DashMember, error) {
	return findOneByField("id", id)
}

func GetAllByDashId(id int) (types.DashMembers, error) {
	return findAllByField("dash_id", id, 0)
}

func GetAllByUserId(id int) (types.DashMembers, error) {
	return findAllByField("user_id", id, 0)
}

func Create(m *types.DashMember) error {
	conn := mysql.Conn()

	values := []interface{}{m.DashId, m.UserId, types.MemberStatusApproved}
	sqltext := "INSERT INTO dashboard_members (dash_id, user_id, status) VALUES (?, ?, ?)"

	res, err := conn.Exec(sqltext, values...)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = int(id)

	return nil
}

func Remove(id int) error {
	conn := mysql.Conn()

	_, err := conn.Exec("DELETE FROM dashboard_members WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
