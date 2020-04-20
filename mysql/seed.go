package mysql

import (
	"github.com/504dev/kidlog/cipher"
	"github.com/504dev/kidlog/types"
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
	values := []interface{}{1, 1, "System"}
	sql := "INSERT INTO dashboards (id, owner_id, name) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE name=name"
	_, err := db.Exec(sql, values...)
	if err != nil {
		panic(err)
	}
}

func SeedKeys() {
	pub, priv, _ := cipher.GenerateKeyPairBase64(32)
	values := []interface{}{1, 1, "Default", pub, priv}
	sql := "INSERT INTO dashboard_keys (id, dash_id, name, public_key, private_key) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=name"
	_, err := db.Exec(sql, values...)
	if err != nil {
		panic(err)
	}
}
