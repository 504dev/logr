package mysql

import (
	"fmt"
	"github.com/504dev/kidlog/config"
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
	Schema()
}

func Schema() {
	var err error
	tables := []string{"dashboards", "dashboard_members", "users"}
	for _, table := range tables {
		path := fmt.Sprintf("../../mysql/schema/%v.sql", table)
		sql, _ := ioutil.ReadFile(path)
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}
}
