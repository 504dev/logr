package logger

import (
	"encoding/json"
	"fmt"
	Log "github.com/504dev/kidlog/models/log"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

type Logger struct {
}

var commit = readCommit()
var pid = os.Getpid()

func caller() string {
	stack := string(debug.Stack())
	caller := strings.TrimSpace(strings.Split(stack, "\n")[10])
	splitted := regexp.MustCompile(`[\s\/]+`).Split(caller, 20)
	length := len(splitted)
	caller = strings.Join(splitted[length-3:length-1], "/")
	return caller
}

func prefix(level string) string {
	dt := time.Now().Format(time.RFC3339)
	return fmt.Sprintf("[KID] %v %v [%v, pid=%v, %v]", dt, level, commit[:6], pid, caller())
}

func readCommit() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	stdout, err := cmd.Output()

	if err != nil {
		log.Println(err.Error)
		return ""
	}

	return string(stdout)
}

func Info(v ...interface{}) {
	args := []interface{}{prefix("info")}
	args = append(args, v...)
	fmt.Println(args...)
}

func Error(v ...interface{}) {
	args := []interface{}{prefix("info")}
	args = append(args, v...)
	log.Println(args...)
}

type Handler struct{}

var publicKey = "MDwwDQYJKoZIhvcNAQEBBQADKwAwKAIhAJ/dD/71Ak7I0e0Q9R9M+yQr45dk3FxbSLT6PQ3vPwLLAgMBAAE="
var privateKey = "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q=="

func (t Handler) Write(b []byte) (int, error) {
	conn, err := net.Dial("udp", "127.0.0.1:7776")
	if err != nil {
		return 0, err
	}
	logitem := Log.Log{
		DashId:    0,
		Timestamp: time.Now().UnixNano(),
		Hostname:  "kidlog.loc",
		Logname:   "gin",
		Level:     "info",
		Message:   string(b),
	}
	cipherText, err := logitem.Encrypt(privateKey)
	if err != nil {
		return 0, err
	}
	lp := Log.LogPackage{
		PublicKey:  publicKey,
		CipherText: cipherText,
	}
	msg, err := json.Marshal(lp)
	if err != nil {
		return 0, err
	}
	_, err = conn.Write(msg)
	if err != nil {
		return 0, err
	}
	conn.Close()
	return len(b), nil
}
