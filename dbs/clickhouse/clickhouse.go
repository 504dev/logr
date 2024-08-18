package clickhouse

import (
	"fmt"
	"github.com/504dev/logr/config"
	_ "github.com/ClickHouse/clickhouse-go"
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
	db, err = sqlx.Connect("clickhouse", config.Get().Clickhouse)

	if err == nil {
		Migrate()
		return
	}

	if retries > 0 {
		fmt.Fprintf(os.Stderr, "(%v) clickhouse connect retry: %s\n", retries, err)
		<-time.After(time.Second)
		MustInit(retries - 1)
		return
	}

	panic(err)
}

func Migrate() {
	db, _ := sqlx.Connect("clickhouse", config.Get().Clickhouse+"&x-multi-statement=true")
	_ = goose.SetDialect("clickhouse")
	err := goose.Up(db.DB, "./dbs/clickhouse/migrations")
	if err != nil {
		panic(err)
	}
}
