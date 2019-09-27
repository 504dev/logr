package clickhouse

import (
	"github.com/504dev/kidlog/types"
)

func Init() {
	// create connection sqlx
	Schema()
}

func Schema() {
	// create schemas if not exist
}

func InsertQueue(log types.Log) {

}

func InsertBatch(logs types.Logs) {

}
