package clickhouse

import (
	"fmt"
	"github.com/504dev/kidlog/config"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
)

var db *sqlx.DB

func Conn() *sqlx.DB {
	return db
}

func Init() {
	var err error
	db, err = sqlx.Open("clickhouse", config.Get().Clickhouse)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	Schema()
}

func Schema() {
	var err error
	tables := []string{"logs"}
	for _, table := range tables {
		path := fmt.Sprintf("../../clickhouse/schema/%v.sql", table)
		sql, _ := ioutil.ReadFile(path)
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}
}
