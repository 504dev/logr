package mysql

import (
	"fmt"
	"os"
)

func Schema() {
	var err error
	tables := []string{"users", "dashboards", "dashboard_members", "dashboard_keys"}
	for _, table := range tables {
		path := fmt.Sprintf("./mysql/schema/%v.sql", table)
		sql, _ := os.ReadFile(path)
		_, err = db.Exec(string(sql))
		if err != nil {
			panic(err)
		}
	}
}
