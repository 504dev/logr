package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
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

func Create(ownerId int, name string) (*types.Dashboard, error) {
	conn := mysql.Conn()

	pubkey, privkey, err := cipher.GenerateKeyPairBase64(256)
	if err != nil {
		return nil, err
	}

	values := []interface{}{ownerId, name, pubkey, privkey}
	sql := "INSERT INTO dashboards (owner_id, name, public_key, private_key) VALUES (?, ?, ?, ?)"

	res, err := conn.Exec(sql, values...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	dashboard, err := GetById(int(id))

	return dashboard, err
}

func Share(m *types.DashMember) error {
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
