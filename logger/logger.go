package logger

import (
	"github.com/504dev/go-kidlog"
)

var conf = go_kidlog.Config{
	Udp:        "127.0.0.1:7776",
	PublicKey:  "MDwwDQYJKoZIhvcNAQEBBQADKwAwKAIhAJ/dD/71Ak7I0e0Q9R9M+yQr45dk3FxbSLT6PQ3vPwLLAgMBAAE=",
	PrivateKey: "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q==",
	Hostname:   "localhost",
}

var Logr, _ = conf.Create("main.log")

func Info(v ...interface{}) {
	Logr.Info(v...)
}

func Error(v ...interface{}) {
	Logr.Error(v...)
}
