package clickhouse

import (
	"github.com/504dev/kidlog/models/log"
)

func Init() {
	// create connection sqlx
	Schema()
}

func Schema() {
	// create schemas if not exist
}

func InsertQueue(log log.Log) {

}

func InsertBatch(logs log.Logs) {

}
