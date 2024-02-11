package mysql

import (
	"fmt"
	"github.com/504dev/logr/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

var db *sqlx.DB

func Conn() *sqlx.DB {
	return db
}

func Init(retries int) {
	var err error
	db, err = sqlx.Open("mysql", config.Get().Mysql+"?parseTime=true")
	if err == nil {
		err = db.Ping()
		if err == nil {
			Schema()
			SeedUsers()
			SeedDashboards()
			SeedKeys()
			return
		}
	}
	if retries > 0 {
		fmt.Fprintf(os.Stderr, "(%v) mysql connect retry: %s\n", retries, err)
		<-time.After(time.Second)
		Init(retries - 1)
		return
	}
	panic(err)
}
