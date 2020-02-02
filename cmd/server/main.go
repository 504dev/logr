package main

import (
	"github.com/504dev/kidlog/clickhouse"
	"github.com/504dev/kidlog/config"
	"github.com/504dev/kidlog/mysql"
	"github.com/504dev/kidlog/server"
)

func main() {
	config.Init()
	clickhouse.Init()
	mysql.Init()
	go server.Udp()
	server.Init()
}
