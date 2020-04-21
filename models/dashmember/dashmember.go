package dashmember

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
)

func findAllByField(fieldname string, val interface{}, limit int) (types.DashMembers, error) {
	conn := mysql.Conn()
	members := types.DashMembers{}
	sqltext := fmt.Sprintf("SELECT id, dash_id, user_id FROM dashboard_members WHERE %v = ?", fieldname)
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

func GetByDashId(id int) (types.DashMembers, error) {
	return findAllByField("dash_id", id, 0)
}
