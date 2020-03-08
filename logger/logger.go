package logger

import (
	"github.com/504dev/go-kidlog"
	"os"
	"strconv"
)

var hostname, _ = os.Hostname()

var conf = go_kidlog.Config{
	Udp:    "127.0.0.1:7776",
	DashId: 1,
	//PublicKey:  "MDwwDQYJKoZIhvcNAQEBBQADKwAwKAIhAJ/dD/71Ak7I0e0Q9R9M+yQr45dk3FxbSLT6PQ3vPwLLAgMBAAE=",
	PrivateKey: "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q==",
	Hostname:   hostname,
}

var Main, _ = conf.Create("main.log")
var gin, _ = conf.Create("gin.log")
var Gin = gin.Parser(func(log *go_kidlog.Log) {
	codestr := log.Message[38:41]
	code, _ := strconv.Atoi(codestr)
	if code >= 400 && code <= 499 {
		log.Level = go_kidlog.LevelWarn
	} else if code >= 500 && code <= 599 {
		log.Level = go_kidlog.LevelError
	}
})

var Create = conf.Create
var Info = Main.Info
var Error = Main.Error
var Debug = Main.Debug
var Warn = Main.Warn
