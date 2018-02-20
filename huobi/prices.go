package huobi

import (
	"log"
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
	priceDetail := GetMarketDetail(sym)
	if priceDetail.ErrMsg != "" {
		log.Print("[-] ErrMsg from huobi: ", priceDetail.ErrMsg)
	}
	price := priceDetail.Tick.Close
	return strconv.FormatFloat(price, 'f', -1, 32)
}
