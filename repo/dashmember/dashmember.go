package dashmember

import (
	"fmt"
	"github.com/504dev/logr/dbs/mysql"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
)

type DashboardMemberRepo struct {
	conn *sqlx.DB
}

func NewDashboardMemberRepo() *DashboardMemberRepo {
	return &DashboardMemberRepo{
		conn: mysql.Conn(),
	}
}

func (repo *DashboardMemberRepo) findAllByField(
	fieldname string,
	val interface{},
	limit int,
) (types.DashMembers, error) {
	members := types.DashMembers{}
	sqltext := fmt.Sprintf("SELECT id, dash_id, user_id, status FROM dashboard_members WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	err := repo.conn.Select(&members, sqltext, val)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (repo *DashboardMemberRepo) findOneByField(fieldname string, val interface{}) (*types.DashMember, error) {
	members, err := repo.findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, nil
	}
	return members[0], nil
}

func (repo *DashboardMemberRepo) GetById(id int) (*types.DashMember, error) {
	return repo.findOneByField("id", id)
}

func (repo *DashboardMemberRepo) GetByDashId(id int) (types.DashMembers, error) {
	return repo.findAllByField("dash_id", id, 0)
}

func (repo *DashboardMemberRepo) GetByUserId(id int) (types.DashMembers, error) {
	return repo.findAllByField("user_id", id, 0)
}

func (repo *DashboardMemberRepo) Create(m *types.DashMember) error {
	values := []interface{}{m.DashId, m.UserId, types.MEMBER_STATUS_APPROVED}
	sqltext := "INSERT INTO dashboard_members (dash_id, user_id, status) VALUES (?, ?, ?)"

	res, err := repo.conn.Exec(sqltext, values...)

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

func (repo *DashboardMemberRepo) Remove(id int) error {
	_, err := repo.conn.Exec("DELETE FROM dashboard_members WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
