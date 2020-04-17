package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
	"strings"
)

func GetAll() (types.Dashboards, error) {
	conn := mysql.Conn()
	dashboards := types.Dashboards{}
	err := conn.Select(&dashboards, "SELECT id, owner_id, name, public_key, private_key FROM dashboards")
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func findAllByField(fieldname string, val interface{}, limit int) (types.Dashboards, error) {
	conn := mysql.Conn()
	dashboards := types.Dashboards{}
	sql := fmt.Sprintf("SELECT id, owner_id, name, public_key, private_key FROM dashboards WHERE %v = ?", fieldname)
	if limit > 0 {
		sql = fmt.Sprintf("%v LIMIT %v", sql, limit)
	}
	err := conn.Select(&dashboards, sql, val)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func findOneByField(fieldname string, val interface{}) (*types.Dashboard, error) {
	dashboards, err := findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(dashboards) == 0 {
		return nil, nil
	}
	return dashboards[0], nil
}

func GetById(id int) (*types.Dashboard, error) {
	return findOneByField("id", id)
}

func GetByPub(pub string) (*types.Dashboard, error) {
	return findOneByField("public_key", pub)
}

func GetUserDashboards(id int) (types.Dashboards, error) {
	return findAllByField("owner_id", id, 0)
}

func GetShared(id int) (types.Dashboards, error) {
	conn := mysql.Conn()
	members := types.DashMembers{}
	sql := "SELECT id, dash_id, user_id FROM dashboard_members WHERE user_id = ?"
	err := conn.Select(&members, sql, id)
	if err != nil {
		return nil, err
	}
	ids := append(members.DashIds(), 1)
	dashboards := types.Dashboards{}
	if len(ids) == 0 {
		return dashboards, nil
	}
	sql = fmt.Sprintf("SELECT id, owner_id, name FROM dashboards WHERE id IN (%v)", strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]"))
	err = conn.Select(&dashboards, sql)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func Create(dash *types.Dashboard) error {
	conn := mysql.Conn()

	var err error
	dash.PublicKey, dash.PrivateKey, err = cipher.GenerateKeyPairBase64(256)
	if err != nil {
		return err
	}

	values := []interface{}{dash.OwnerId, dash.Name, dash.PublicKey, dash.PrivateKey}
	sql := "INSERT INTO dashboards (owner_id, name, public_key, private_key) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=name"

	res, err := conn.Exec(sql, values...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	dash.Id = int(id)

	return nil
}

func Update(dash *types.Dashboard) error {
	conn := mysql.Conn()

	values := []interface{}{dash.Name, dash.Id}
	sql := "UPDATE dashboards SET name = ? where id = ?"

	_, err := conn.Exec(sql, values...)
	if err != nil {
		return err
	}

	item, err := GetById(dash.Id)
	if err != nil {
		return err
	}

	*dash = *item

	return nil
}

func Delete(id int) error {
	conn := mysql.Conn()

	_, err := conn.Exec("DELETE FROM dashboards where id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func AddMember(m *types.DashMember) error {
	conn := mysql.Conn()

	values := []interface{}{m.DashId, m.UserId}
	sql := "INSERT INTO dashboard_members (dash_id, user_id) VALUES (?, ?)"

	res, err := conn.Exec(sql, values...)

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
