package huobi

import (
	huobiAPI "github.com/jm33-m0/huobi-go/services"
	"strconv"
)

// GetPrices : Get a list of coin prices
func GetPrices() string {
	btcUSDT := getClosePrice("btcusdt")
	ethBTC := getClosePrice("ethbtc")
	htBTC := getClosePrice("htbtc")
	htUSDT := getClosePrice("htusdt")

	retVal := "`BTC - USDT : " + btcUSDT + "\n" + "ETH - BTC : " + ethBTC + "\n" + "HT - BTC : " + htBTC + "\n" + "HT - USDT : " + htUSDT + "`"

	return retVal
}

func getClosePrice(sym string) string {
	return strconv.FormatFloat(huobiAPI.GetMarketDetail(sym).Tick.Close, 'f', -1, 32)
}
