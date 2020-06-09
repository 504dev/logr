package mysql

import (
	"github.com/504dev/logr/cipher"
	"github.com/504dev/logr/types"
)

func SeedUsers() {
	values := []interface{}{1, 0, "admin", types.RoleAdmin}
	sql := "INSERT INTO users (id, github_id, username, role) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE username=username"
	_, err := db.Exec(sql, values...)
	if err != nil {
		panic(err)
	}
}
func SeedDashboards() {
	sqltext := "INSERT INTO dashboards (id, owner_id, name) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE name=name"
	stmt, err := db.Prepare(sqltext)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	batch := [][]interface{}{
		{types.DashboardSystemId, 1, "System"},
		{types.DashboardDemoId, 1, "Demo"},
	}
	for _, v := range batch {
		_, err = stmt.Exec(v...)
		if err != nil {
			panic(err)
		}
	}
}

func SeedKeys() {
	sqltext := "INSERT INTO dashboard_keys (id, dash_id, name, public_key, private_key) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=name"
	stmt, err := db.Prepare(sqltext)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	batch := [][]interface{}{
		{1, types.DashboardSystemId, "Default"},
		{2, types.DashboardDemoId, "Default"},
	}
	for _, v := range batch {
		pub, priv, _ := cipher.GenerateKeyPairBase64(32)
		v = append(v, pub, priv)
		_, err = stmt.Exec(v...)
		if err != nil {
			panic(err)
		}
	}
}
