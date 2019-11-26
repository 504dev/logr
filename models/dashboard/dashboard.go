package dashboard

import (
	"fmt"
	"github.com/504dev/kidlog/mysql"
)

var dashboards = Dashboards{
	{
		Id:         1,
		OwnerId:    1,
		Name:       "sandbox",
		PublicKey:  "MDwwDQYJKoZIhvcNAQEBBQADKwAwKAIhAJ/dD/71Ak7I0e0Q9R9M+yQr45dk3FxbSLT6PQ3vPwLLAgMBAAE=",
		PrivateKey: "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q==",
	},
}

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

func getByField(fieldname string, val interface{}) (*Dashboard, error) {
	var dash Dashboard
	conn := mysql.Conn()
	sql := fmt.Sprintf("SELECT id, owner_id, name, public_key, private_key FROM dashboards WHERE %v = ?", fieldname)
	row := conn.QueryRowx(sql, val)
	err := row.StructScan(&dash)

	if err != nil {
		return nil, err
	}

	return &dash, nil
}

func GetById(id int) (*Dashboard, error) {
	return getByField("id", id)
}

func GetByPub(pub string) (*Dashboard, error) {
	return getByField("public_key", pub)
}
