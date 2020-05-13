package logger

import (
	"encoding/json"
	"fmt"
	logr "github.com/504dev/logr-go-client"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"time"
)

func crypto(conf *logr.Config) {
	l, _ := conf.NewLogger("crypto.log")
	for {
		time.Sleep(60 * time.Second)
		day := time.Now().Format("2006-01-02")
		path := fmt.Sprintf("/get-day-snapshot?day=%v&uni=1&format=ohlcv", day)
		bytes, err := request(path)
		if err != nil {
			l.Error(err)
			continue
		}
		prices := map[string]map[string]map[string]float64{}

		if err = json.Unmarshal(bytes, &prices); err != nil {
			l.Error(err)
			continue
		}

		for _, base := range [3]string{"BTC", "ETH", "LTC"} {
			sym := base + "_USDT"
			l.Info(
				"%v price: %v %v$, %v %v$, %v %v$",
				color.New(color.Bold).SprintFunc()(base),
				color.CyanString("HitBTC"),
				prices["hitbtc"][sym]["c"],
				color.YellowString("Binance"),
				prices["binance"][sym]["c"],
				color.GreenString("Bitfinex"),
				prices["bitfinex"][sym]["c"],
			)
			l.Info(
				"%v volume: %v %.0f$, %v %.0f$, %v %.0f$",
				color.New(color.Bold).SprintFunc()(base),
				color.CyanString("HitBTC"),
				prices["hitbtc"][sym]["v"],
				color.YellowString("Binance"),
				prices["binance"][sym]["v"],
				color.GreenString("Bitfinex"),
				prices["bitfinex"][sym]["v"],
			)
		}

		l.Debug(string(bytes))
	}
}

func request(path string) ([]byte, error) {
	res := []byte{}
	url := fmt.Sprintf("http://212.224.113.196:5554%v", path)
	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return bodyBytes, nil
}
