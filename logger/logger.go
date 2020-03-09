package logger

import (
	. "github.com/504dev/go-kidlog"
	"os"
	"strconv"
)

var hostname, _ = os.Hostname()

var conf = Config{
	Udp:        "127.0.0.1:7776",
	DashId:     1,
	PrivateKey: "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q==",
	Hostname:   hostname,
}

var Main, _ = conf.Create("main.log")
var gin, _ = conf.Create("gin.log")
var Gin = gin.CustomWritter(func(log *Log) {
	codestr := log.Message[38:41]
	code, _ := strconv.Atoi(codestr)
	if code >= 400 && code <= 499 {
		log.Level = LevelWarn
	} else if code >= 500 && code <= 599 {
		log.Level = LevelError
	}
})

var Create = conf.Create
var Info = Main.Info
var Error = Main.Error
var Debug = Main.Debug
var Warn = Main.Warn
