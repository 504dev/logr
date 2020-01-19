package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/mysql"
)

func GetAll() (Dashboards, error) {
	conn := mysql.Conn()
	rows, err := conn.Queryx("SELECT id, owner_id, name, public_key, private_key FROM dashboards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dashboards := make(Dashboards, 0)

	for rows.Next() {
		var dash Dashboard
		err := rows.StructScan(&dash)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, &dash)
	}
	return dashboards, nil
}

func getAllByField(fieldname string, val interface{}) (Dashboards, error) {
	var dashboards Dashboards
	conn := mysql.Conn()
	sql := fmt.Sprintf("SELECT id, owner_id, name, public_key, private_key FROM dashboards WHERE %v = ?", fieldname)
	rows, err := conn.Queryx(sql, val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dashboards = make(Dashboards, 0)
	for rows.Next() {
		dash := Dashboard{}
		err = rows.StructScan(&dash)
		dashboards = append(dashboards, &dash)

	}

	return dashboards, nil
}

func getOneByField(fieldname string, val interface{}) (*Dashboard, error) {
	dashboards, err := getAllByField(fieldname, val)
	if err != nil {
		return nil, err
	}
	if len(dashboards) == 0 {
		return nil, nil
	}
	return dashboards[0], nil
}

func GetById(id int) (*Dashboard, error) {
	return getOneByField("id", id)
}

func GetByPub(pub string) (*Dashboard, error) {
	return getOneByField("public_key", pub)
}

func GetUserDashboards(id int) (Dashboards, error) {
	return getAllByField("owner_id", id)
}

func CreateDashboard(ownerId int, name string) (*Dashboard, error) {
	conn := mysql.Conn()

	pubkey, privkey, err := cipher.GenerateKeyPairBase64(256)
	if err != nil {
		return nil, err
	}

	values := []interface{}{ownerId, name, pubkey, privkey}
	sqlstr := `INSERT INTO dashboards (owner_id, name, public_key, private_key) VALUES (?, ?, ?, ?)`
	stmt, err := conn.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user, err := GetById(int(id))

	return user, err
}
