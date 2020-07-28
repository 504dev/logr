package logger

import (
	"encoding/json"
	"fmt"
	logr "github.com/504dev/logr-go-client"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"time"
)

func crypto(conf *logr.Config) {
	l, _ := conf.NewLogger("crypto.log")
	for {
		time.Sleep(30 * time.Second)
		delta := l.Time("pricer:/get-day-snapshot", time.Millisecond)
		day := time.Now().Format("2006-01-02")
		path := fmt.Sprintf("/get-day-snapshot?day=%v&uni=1&format=ohlcv&quote=USDT", day)
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
		delta()

		for _, base := range [3]string{"BTC", "ETH", "LTC"} {
			sym := base + "_USDT"
			hitp := prices["hitbtc"][sym]["c"]
			binp := prices["binance"][sym]["c"]
			bitp := prices["bitfinex"][sym]["c"]
			l.Info(
				"%v price: %v %v$, %v %v$, %v %v$",
				color.New(color.Bold).SprintFunc()(base),
				color.CyanString("HitBTC"),
				humanize.Commaf(hitp),
				color.YellowString("Binance"),
				humanize.Commaf(binp),
				color.GreenString("Bitfinex"),
				humanize.Commaf(bitp),
			)
			l.Touch(fmt.Sprintf("price:%v", sym)).Avg(hitp).Avg(binp).Avg(bitp).Min(hitp).Min(binp).Min(bitp).Max(hitp).Max(binp).Max(bitp)
			hitv := prices["hitbtc"][sym]["v"]
			binv := prices["binance"][sym]["v"]
			bitv := prices["bitfinex"][sym]["v"]
			l.Info(
				"%v volume: %v %v$, %v %v$, %v %v$",
				color.New(color.Bold).SprintFunc()(base),
				color.CyanString("HitBTC"),
				humanize.Comma(int64(hitv)),
				color.YellowString("Binance"),
				humanize.Comma(int64(binv)),
				color.GreenString("Bitfinex"),
				humanize.Comma(int64(bitv)),
			)
			l.Avg(fmt.Sprintf("volume:%v", sym), hitp+binv+bitv)
			l.Info(
				"%v price widget %v",
				color.New(color.Bold).SprintFunc()(base),
				l.Widget("max", fmt.Sprintf("price:%v", sym), 30),
			)
			if sym == "BTC_USDT" {
				v := hitv + bitv + binv
				l.Per("volume:BTC_USDT:hitbtc", hitv, v)
				l.Per("volume:BTC_USDT:bitfinex", bitv, v)
				l.Per("volume:BTC_USDT:binance", binv, v)
			}
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
