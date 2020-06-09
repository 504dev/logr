package dashkey

import (
	"database/sql"
	"fmt"
	"github.com/504dev/logr/cipher"
	"github.com/504dev/logr/mysql"
	"github.com/504dev/logr/types"
)

func Create(key *types.DashKey, tx *sql.Tx) error {

	var err error
	var values []interface{}
	var sqltext string
	var res sql.Result

	key.PublicKey, key.PrivateKey, err = cipher.GenerateKeyPairBase64(256)
	if err != nil {
		return err
	}

	values = []interface{}{key.DashId, key.Name, key.PublicKey, key.PrivateKey}
	sqltext = "INSERT INTO dashboard_keys (dash_id, name, public_key, private_key) VALUES (?, ?, ?, ?)"

	if tx != nil {
		res, err = tx.Exec(sqltext, values...)
	} else {
		res, err = mysql.Conn().Exec(sqltext, values...)
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

func findAllByField(fieldname string, val interface{}, limit int) (types.DashKeys, error) {
	conn := mysql.Conn()
	keys := types.DashKeys{}
	sqltext := fmt.Sprintf("SELECT id, dash_id, name, public_key, private_key FROM dashboard_keys WHERE %v = ?", fieldname)
	if limit > 0 {
		sqltext = fmt.Sprintf("%v LIMIT %v", sqltext, limit)
	}
	err := conn.Select(&keys, sqltext, val)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func findOneByField(fieldname string, val interface{}) (*types.DashKey, error) {
	keys, err := findAllByField(fieldname, val, 1)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}
	return keys[0], nil
}

func GetById(id int) (*types.DashKey, error) {
	return findOneByField("id", id)
}

func GetByDashId(id int) (types.DashKeys, error) {
	return findAllByField("dash_id", id, 0)
}

func GetByPub(pub string) (*types.DashKey, error) {
	return findOneByField("public_key", pub)
}
