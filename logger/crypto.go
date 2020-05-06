package logger

import (
	"fmt"
	logr "github.com/504dev/logr-go-client"
	"io/ioutil"
	"net/http"
	"time"
)

func crypto(conf *logr.Config) {
	l, _ := conf.NewLogger("crypto.log")
	for {
		day := time.Now().Format("2006-01-02")
		path := fmt.Sprintf("/get-day-snapshot?day=%v&uni=1", day)
		j, err := request(path)
		if err != nil {
			l.Error(err)
			continue
		}
		l.Debug(j)
		time.Sleep(10 * time.Second)
	}
}

func request(path string) (string, error) {
	url := fmt.Sprintf("http://212.224.113.196:5554%v", path)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
