package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/models/dashkey"
	"github.com/504dev/kidlog/models/dashmember"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/types"
	"strings"
)

func GetAll() (types.Dashboards, error) {
	conn := mysql.Conn()
	dashboards := types.Dashboards{}
	err := conn.Select(&dashboards, "SELECT id, owner_id, name FROM dashboards")
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func findAllByField(fieldname string, val interface{}, limit int) (types.Dashboards, error) {
	conn := mysql.Conn()
	dashboards := types.Dashboards{}
	sqltext := fmt.Sprintf("SELECT id, owner_id, name FROM dashboards WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	err := conn.Select(&dashboards, sqltext, val)
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

func GetUserDashboards(id int) (types.Dashboards, error) {
	return findAllByField("owner_id", id, 0)
}

func GetSystemIds(role int) []int {
	ids := make([]int, 0, 2)
	if role != 0 {
		ids = append(ids, types.DashboardDemoId)
	}
	if role == types.RoleAdmin {
		ids = append(ids, types.DashboardSystemId)
	}
	return ids
}
func GetShared(id int, role int) (types.Dashboards, error) {
	members, err := dashmember.GetAllByUserId(id)
	if err != nil {
		return nil, err
	}
	ids := members.DashIds()
	ids = append(ids, GetSystemIds(role)...)
	if len(ids) == 0 {
		return types.Dashboards{}, nil
	}
	return GetByIds(ids)
}

func GetByIds(ids []int) (types.Dashboards, error) {
	conn := mysql.Conn()
	dashboards := types.Dashboards{}
	inString := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	sqltext := fmt.Sprintf("SELECT id, owner_id, name FROM dashboards WHERE id IN (%v)", inString)
	err := conn.Select(&dashboards, sqltext)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func Create(dash *types.Dashboard) error {
	conn := mysql.Conn()

	var err error
	var values []interface{}
	var sqltext string

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	values = []interface{}{dash.OwnerId, dash.Name}
	sqltext = "INSERT INTO dashboards (owner_id, name) VALUES (?, ?)"

	res, err := tx.Exec(sqltext, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	dashId := int(id)

	dk := &types.DashKey{DashId: dashId, Name: "Default"}
	err = dashkey.Create(dk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	dash.Id = dashId
	dash.Keys = types.DashKeys{dk}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func Update(dash *types.Dashboard) error {
	conn := mysql.Conn()

	values := []interface{}{dash.Name, dash.Id}
	sqltext := "UPDATE dashboards SET name = ? where id = ?"

	_, err := conn.Exec(sqltext, values...)
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
