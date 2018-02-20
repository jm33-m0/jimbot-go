package huobi

import (
	"log"
	"strconv"
)

// GetPrices : Get a list of coin prices
func GetPrices() string {
	btcUSDT := getClosePrice("btcusdt")
	ethUSDT := getClosePrice("ethusdt")
	ethBTC := getClosePrice("ethbtc")
	htUSDT := getClosePrice("htusdt")
	htBTC := getClosePrice("htbtc")

	retVal := "`BTC - USDT : " + btcUSDT + "\nETH - USDT : " + ethUSDT + "\nETH - BTC  : " + ethBTC + "\nHT  - USDT : " + htUSDT + "\nHT  - BTC  : " + htBTC + "`"

	return retVal
}

func getClosePrice(sym string) string {
	priceDetail := GetMarketDetail(sym)
	if priceDetail.ErrMsg != "" {
		log.Print("[-] ErrMsg from huobi: ", priceDetail.ErrMsg)
	}
	price := priceDetail.Tick.Close
	return strconv.FormatFloat(price, 'f', -1, 32)
}
