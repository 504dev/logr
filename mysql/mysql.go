package mysql

import (
	"github.com/504dev/kidlog/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

var db *sqlx.DB

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
	users, _ := ioutil.ReadFile("../../mysql/schema/users.sql")
	_, err = db.Exec(string(users))
	if err != nil {
		panic(err)
	}
	dashboards, _ := ioutil.ReadFile("../../mysql/schema/dashboards.sql")
	_, err = db.Exec(string(dashboards))
	if err != nil {
		panic(err)
	}
}
