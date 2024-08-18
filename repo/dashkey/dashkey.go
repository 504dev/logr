package dashkey

import (
	"database/sql"
	"fmt"
	"github.com/504dev/logr/dbs/mysql"
	"github.com/504dev/logr/libs/cipher"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

type DashboardKeyRepo struct {
	conn *sqlx.DB
}

func NewDashboardKeyRepo() *DashboardKeyRepo {
	return &DashboardKeyRepo{
		conn: mysql.Conn(),
	}
}

func (repo *DashboardKeyRepo) Create(key *types.DashKey, tx *sql.Tx) error {
	var err error
	var values []interface{}
	var sqltext string
	var res sql.Result

	const keyBitSize = 256
	key.PublicKey, key.PrivateKey, err = cipher.GenerateKeyPairBase64(keyBitSize)
	if err != nil {
		return err
	}

	values = []interface{}{key.DashId, key.Name, key.PublicKey, key.PrivateKey}
	sqltext = "INSERT INTO dashboard_keys (dash_id, name, public_key, private_key) VALUES (?, ?, ?, ?)"

	if tx != nil {
		res, err = tx.Exec(sqltext, values...)
	} else {
		res, err = repo.conn.Exec(sqltext, values...)
	}

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	key.Id = int(id)

	return nil
}

func (repo *DashboardKeyRepo) findAllByField(fieldname string, val interface{}, limit int) (types.DashKeys, error) {
	sqltext := fmt.Sprintf("SELECT id, dash_id, name, public_key, private_key FROM dashboard_keys WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	keys := types.DashKeys{}
	err := repo.conn.Select(&keys, sqltext, val)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (repo *DashboardKeyRepo) findOneByField(fieldname string, val interface{}) (*types.DashKey, error) {
	keys, err := repo.findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}
	return keys[0], nil
}

func (repo *DashboardKeyRepo) GetById(id int) (*types.DashKey, error) {
	return repo.findOneByField("id", id)
}

func (repo *DashboardKeyRepo) GetByDashId(id int) (types.DashKeys, error) {
	return repo.findAllByField("dash_id", id, 0)
}

func (repo *DashboardKeyRepo) GetByPub(pub string) (*types.DashKey, error) {
	return repo.findOneByField("public_key", pub)
}

func (repo *DashboardKeyRepo) GetByIds(ids []int) (types.DashKeys, error) {
	placeholders := make([]string, len(ids))
	for i, v := range ids {
		placeholders[i] = strconv.Itoa(v)
	}
	sqltext := fmt.Sprintf(
		"SELECT id, dash_id, name, public_key, private_key, name FROM dashboard_keys WHERE id IN (%s)",
		strings.Join(placeholders, ","),
	)
	dashkeys := types.DashKeys{}
	err := repo.conn.Select(&dashkeys, sqltext)
	if err != nil {
		return nil, err
	}
	return dashkeys, nil
}
