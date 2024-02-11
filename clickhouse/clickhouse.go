package clickhouse

import (
	"fmt"
	"github.com/504dev/logr/config"
	_ "github.com/ClickHouse/clickhouse-go"
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
	db, err = sqlx.Open("clickhouse", config.Get().Clickhouse)
	if err == nil {
		err = db.Ping()
		if err == nil {
			Schema()
			return
		}
	}
	if retries > 0 {
		fmt.Fprintf(os.Stderr, "(%v) clickhouse connect retry: %s\n", retries, err)
		<-time.After(time.Second)
		Init(retries - 1)
		return
	}
	panic(err)
}

func Schema() {
	var err error
	tables := []string{"logs", "counts"}
	for _, table := range tables {
		path := fmt.Sprintf("./clickhouse/schema/%v.sql", table)
		sql, _ := os.ReadFile(path)
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}
}
