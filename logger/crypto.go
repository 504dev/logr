package logger

import (
	"encoding/json"
	"fmt"
	logr "github.com/504dev/logr-go-client"
	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"io"
	"net/http"
	"time"
)

type BinancePrice struct {
	Price  float64 `json:"lastPrice,string"`
	Volume float64 `json:"quoteVolume,string"`
}

type HitbtcPrice struct {
	Price  float64 `json:"last,string"`
	Volume float64 `json:"volume_quote,string"`
}

type BitfinexPrice struct {
	Price  float64
	Volume float64
}

func (b *BitfinexPrice) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{nil, nil, nil, nil, nil, nil, &b.Price, &b.Volume}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	b.Volume *= b.Price
	return nil
}

func crypto(conf *logr.Config) {
	l, _ := conf.NewLogger("crypto.log")
	l.Body = "[{version}, pid={pid}] {message}"
	for {
		l.Info("")
		time.Sleep(30 * time.Second)
		l.Info("**************************************************")
		totals := make(map[string]float64)
		for _, base := range [3]string{"BTC", "ETH", "LTC"} {
			l.Info("")
			sym := base + "_USDT"
			bin, hit, bit := BinancePrice{}, HitbtcPrice{}, BitfinexPrice{}
			var err error
			err = request(&bin, fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%vUSDT", base))
			if err != nil {
				Logger.Error("Demo crypto.log binance: %v", err)
				continue
			}
			err = request(&hit, fmt.Sprintf("https://api.hitbtc.com/api/3/public/ticker/%vUSDT", base))
			if err != nil {
				Logger.Error("Demo crypto.log hitbtc: %v", err)
				continue
			}
			err = request(&bit, fmt.Sprintf("https://api-pub.bitfinex.com/v2/ticker/t%vUSD", base))
			if err != nil {
				Logger.Error("Demo crypto.log bitfinex: %v", err)
				continue
			}
			binP, hitP, bitP := bin.Price, hit.Price, bit.Price
			binV, hitV, bitV := bin.Volume, hit.Volume, bit.Volume

			l.Touch(fmt.Sprintf("price:%v", sym)).Avg(hitP).Avg(binP).Avg(bitP)
			l.Avg(fmt.Sprintf("volume:%v", sym), hitV+binV+bitV)

			bold := color.New(color.Bold).SprintFunc()

			l.Info(
				"%v %v %v$ (%v$)",
				color.CyanString("HitBTC"),
				bold(base),
				bold(humanize.Commaf(hitP)),
				humanize.Comma(int64(hitV)),
			)

			l.Info(
				"%v %v %v$ (%v$)",
				color.GreenString("Bitfinex"),
				bold(base),
				bold(humanize.Commaf(bitP)),
				humanize.Comma(int64(bitV)),
			)

			l.Info(
				"%v %v %v$ (%v$)",
				color.HiYellowString("Binance"),
				bold(base),
				bold(humanize.Commaf(binP)),
				humanize.Comma(int64(binV)),
			)

			l.Notice(
				"%v price %v widget!",
				color.New(color.Bold).SprintFunc()(base),
				l.Snippet("max", fmt.Sprintf("price:%v", sym), 30),
			)

			totalV := hitV + bitV + binV
			l.Per("volume:hitbtc", hitV, totalV)
			l.Per("volume:bitfinex", bitV, totalV)
			l.Per("volume:binance", binV, totalV)

			totals[sym] = totalV
			totals[""] += totalV

			time.Sleep(time.Second)
		}

		for sym, val := range totals {
			if sym != "" {
				l.Per(fmt.Sprintf("volume:%v", sym), val, totals[""])
			}
		}
	}
}

func request(res interface{}, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bodyBytes, res); err != nil {
		return err
	}

	return nil
}
