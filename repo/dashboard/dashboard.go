package dashboard

import (
	"fmt"
	"github.com/504dev/logr/config"
	"github.com/504dev/logr/dbs/mysql"
	"github.com/504dev/logr/repo/dashkey"
	"github.com/504dev/logr/repo/dashmember"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DashboardRepo struct {
	conn            *sqlx.DB
	DashboardMember *dashmember.DashboardMemberRepo
	DashboardKey    *dashkey.DashboardKeyRepo
}

func NewDashboardRepo(
	dashboardMemberRepo *dashmember.DashboardMemberRepo,
	dashboardKeyRepo *dashkey.DashboardKeyRepo,
) *DashboardRepo {
	return &DashboardRepo{
		conn:            mysql.Conn(),
		DashboardMember: dashboardMemberRepo,
		DashboardKey:    dashboardKeyRepo,
	}
}

func (repo *DashboardRepo) GetAll() (types.Dashboards, error) {
	dashboards := types.Dashboards{}
	err := repo.conn.Select(&dashboards, "SELECT id, owner_id, name FROM dashboards")
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func (repo *DashboardRepo) findAllByField(fieldname string, val interface{}, limit int) (types.Dashboards, error) {
	dashboards := types.Dashboards{}
	sqltext := fmt.Sprintf("SELECT id, owner_id, name FROM dashboards WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	err := repo.conn.Select(&dashboards, sqltext, val)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func (repo *DashboardRepo) findOneByField(fieldname string, val interface{}) (*types.Dashboard, error) {
	dashboards, err := repo.findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(dashboards) == 0 {
		return nil, nil
	}
	return dashboards[0], nil
}

func (repo *DashboardRepo) GetById(id int) (*types.Dashboard, error) {
	return repo.findOneByField("id", id)
}

func (repo *DashboardRepo) GetByOwnerId(id int) (types.Dashboards, error) {
	return repo.findAllByField("owner_id", id, 0)
}

func (repo *DashboardRepo) GetSystemIds(role int) []int {
	ids := make([]int, 0, 2)
	if role == types.RoleAdmin {
		ids = append(ids, types.DashboardSystemId)
	}
	if config.Get().DemoDash.Enabled {
		ids = append(ids, types.DashboardDemoId)
	}
	return ids
}

func (repo *DashboardRepo) GetByIds(ids []int) (types.Dashboards, error) {
	dashboards := types.Dashboards{}
	inString := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	sqltext := fmt.Sprintf("SELECT id, owner_id, name FROM dashboards WHERE id IN (%v)", inString)
	err := repo.conn.Select(&dashboards, sqltext)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

func (repo *DashboardRepo) GetShared(userId int, role int) (types.Dashboards, error) {
	members, err := repo.DashboardMember.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	ids := members.DashIds()
	ids = append(ids, repo.GetSystemIds(role)...)
	if len(ids) == 0 {
		return types.Dashboards{}, nil
	}
	return repo.GetByIds(ids)
}

func (repo *DashboardRepo) Create(dash *types.Dashboard) error {
	var err error
	var values []interface{}
	var sqltext string

	tx, err := repo.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	values = []interface{}{dash.OwnerId, dash.Name}
	sqltext = "INSERT INTO dashboards (owner_id, name) VALUES (?, ?)"

	res, err := tx.Exec(sqltext, values...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	dashId := int(id)
	dk := &types.DashKey{DashId: dashId, Name: "Default"}
	err = repo.DashboardKey.Create(dk, tx)
	if err != nil {
		return err
	}

	dash.Id = dashId
	dash.Keys = types.DashKeys{dk}

	return tx.Commit()
}

func (repo *DashboardRepo) Update(dash *types.Dashboard) error {
	values := []interface{}{dash.Name, dash.Id}
	sqltext := "UPDATE dashboards SET name = ? where id = ?"

	_, err := repo.conn.Exec(sqltext, values...)
	if err != nil {
		return err
	}

	item, err := repo.GetById(dash.Id)
	if err != nil {
		return err
	}

	*dash = *item

	return nil
}

func (repo *DashboardRepo) Remove(id int) error {
	_, err := repo.conn.Exec("DELETE FROM dashboards where id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
