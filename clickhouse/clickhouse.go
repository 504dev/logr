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

func InsertQueue(log log.Log) bool {
	return false
}

func InsertBatch(logs log.Logs) error {
	return nil
}
