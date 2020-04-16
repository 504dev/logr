package mysql

import (
	"fmt"
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

var db *sqlx.DB

func Conn() *sqlx.DB {
	return db
}

func Init() {
	var err error
	db, err = sqlx.Open("mysql", config.Get().Mysql)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	Schema()
	SeedUsers()
	SeedDashboards()
}

func SeedUsers() {
	values := []interface{}{1, 0, "admin", types.RoleAdmin}
	sql := "INSERT INTO users (id, github_id, username, role) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE username=username"
	_, err := db.Exec(sql, values...)
	if err != nil {
		panic(err)
	}
}
func SeedDashboards() {
	pub, priv, _ := cipher.GenerateKeyPairBase64(32)
	values := []interface{}{1, 1, "System", pub, priv}
	sql := "INSERT INTO dashboards (id, owner_id, name, public_key, private_key) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=name"
	_, err := db.Exec(sql, values...)
	if err != nil {
		panic(err)
	}
}

func Schema() {
	var err error
	tables := []string{"users", "dashboards", "dashboard_members"}
	for _, table := range tables {
		path := fmt.Sprintf("../../mysql/schema/%v.sql", table)
		sql, _ := ioutil.ReadFile(path)
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}
}
