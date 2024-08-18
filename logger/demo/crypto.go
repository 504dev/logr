package demo

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

func crypto(conf *logr.Config, mainlog *logr.Logger) {
	const promptInterval = 30 * time.Second

	cryptolog, _ := conf.NewLogger("crypto.log")
	cryptolog.Body = "[{version}, pid={pid}] {message}"

	for {
		cryptolog.Info("")
		time.Sleep(promptInterval)
		cryptolog.Info("**************************************************")
		totals := make(map[string]float64)

		for _, base := range [3]string{"BTC", "ETH", "LTC"} {
			cryptolog.Info("")
			sym := base + "_USDT"
			bin, hit, bit := BinancePrice{}, HitbtcPrice{}, BitfinexPrice{}
			err := request(&bin, fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%vUSDT", base))
			if err != nil {
				mainlog.Error("Demo crypto.log binance: %v", err)
				continue
			}
			err = request(&hit, fmt.Sprintf("https://api.hitbtc.com/api/3/public/ticker/%vUSDT", base))
			if err != nil {
				mainlog.Error("Demo crypto.log hitbtc: %v", err)
				continue
			}
			err = request(&bit, fmt.Sprintf("https://api-pub.bitfinex.com/v2/ticker/t%vUSD", base))
			if err != nil {
				mainlog.Error("Demo crypto.log bitfinex: %v", err)
				continue
			}
			binP, hitP, bitP := bin.Price, hit.Price, bit.Price
			binV, hitV, bitV := bin.Volume, hit.Volume, bit.Volume

			cryptolog.Touch(fmt.Sprintf("price:%v", sym)).Avg(hitP).Avg(binP).Avg(bitP)
			cryptolog.Avg(fmt.Sprintf("volume:%v", sym), hitV+binV+bitV)

			bold := color.New(color.Bold).SprintFunc()

			cryptolog.Info(
				"%v %v %v$ (%v$)",
				bold(base),
				color.CyanString("HitBTC"),
				bold(humanize.Commaf(hitP)),
				humanize.Comma(int64(hitV)),
			)

			cryptolog.Info(
				"%v %v %v$ (%v$)",
				bold(base),
				color.GreenString("Bitfinex"),
				bold(humanize.Commaf(bitP)),
				humanize.Comma(int64(bitV)),
			)

			cryptolog.Info(
				"%v %v %v$ (%v$)",
				bold(base),
				color.HiYellowString("Binance"),
				bold(humanize.Commaf(binP)),
				humanize.Comma(int64(binV)),
			)

			const snippetSize = 30
			cryptolog.Notice(
				"%v price %v widget!",
				color.New(color.Bold).SprintFunc()(base),
				cryptolog.Snippet("avg", fmt.Sprintf("price:%v", sym), snippetSize),
			)

			totalV := hitV + bitV + binV
			cryptolog.Per("volume:hitbtc", hitV, totalV)
			cryptolog.Per("volume:bitfinex", bitV, totalV)
			cryptolog.Per("volume:binance", binV, totalV)

			totals[sym] = totalV
			totals[""] += totalV

			time.Sleep(time.Second)
		}

		for sym, val := range totals {
			if sym != "" {
				cryptolog.Per(fmt.Sprintf("volume:%v", sym), val, totals[""])
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

	return json.Unmarshal(bodyBytes, res)
}
