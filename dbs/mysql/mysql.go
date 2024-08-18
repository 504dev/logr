package mysql

import (
	"fmt"
	"github.com/504dev/logr/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"os"
	"time"
)

var db *sqlx.DB

func Conn() *sqlx.DB {
	return db
}

func MustInit(retries int) {
	var err error
	db, err = sqlx.Connect("mysql", config.Get().Mysql+"?parseTime=true")

	if err == nil {
		Migrate()
		return
	}

	if retries > 0 {
		fmt.Fprintf(os.Stderr, "(%v) mysql connect retry: %s\n", retries, err)
		<-time.After(time.Second)
		MustInit(retries - 1)
		return
	}

	panic(err)
}

func Migrate() {
	db, _ := sqlx.Connect("mysql", config.Get().Mysql+"?multiStatements=true")
	_ = goose.SetDialect("mysql")
	err := goose.Up(db.DB, "./dbs/mysql/migrations")
	if err != nil {
		panic(err)
	}
}
