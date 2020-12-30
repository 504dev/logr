package mysql

import (
	"github.com/504dev/logr/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Conn() *sqlx.DB {
	return db
}

func Init() {
	var err error
	db, err = sqlx.Open("mysql", config.Get().Mysql+"?parseTime=true")
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
	SeedKeys()
}
