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

var main, _ = conf.NewLogger("main.log")
var gin, _ = conf.NewLogger("gin.log")

var Gin = gin.CustomWritter(func(log *Log) {
	codestr := log.Message[38:41]
	code, _ := strconv.Atoi(codestr)
	if code >= 400 && code <= 499 {
		log.Level = LevelWarn
	} else if code >= 500 && code <= 599 {
		log.Level = LevelError
	}
})

var counter, _ = conf.NewCounter("main.cnt")

var Info = main.Info
var Error = main.Error
var Debug = main.Debug
var Warn = main.Warn

var Inc = counter.Inc
var Max = counter.Max
var Min = counter.Min
var Avg = counter.Avg
var Per = counter.Per
