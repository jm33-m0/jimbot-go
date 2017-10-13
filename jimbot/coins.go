package jimbot

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/buger/jsonparser"
)

// API
const (
	btc = "https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD"
	xmr = "https://min-api.cryptocompare.com/data/price?fsym=XMR&tsyms=BTC,USD"
	eth = "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=BTC,USD"
	etc = "https://min-api.cryptocompare.com/data/price?fsym=ETC&tsyms=BTC,USD"
	bcc = "https://min-api.cryptocompare.com/data/price?fsym=BCH&tsyms=BTC,USD"
)

// CoinPrice : Price info of a coin
type CoinPrice struct {
	CoinName   string
	PriceInBTC string
	PriceInUSD string
}

// GetPrice : Read coin prices from public API
func GetPrice(coin string) CoinPrice {
	var retVal CoinPrice
	var coinURL string

	switch coin {
	case "BTC":
		coinURL = btc
	case "XMR":
		coinURL = xmr
	case "BCC":
		coinURL = bcc
	case "ETH":
		coinURL = eth
	case "ETC":
		coinURL = etc
	default:
		coinURL = btc
	}

	if resp, err := http.Get(coinURL); err == nil {
		defer resp.Body.Close()
		readBody, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			log.Println("[-] Error getting API response")
			log.Print(err)
			return CoinPrice{coin, "n/a", "n/a"}
		}

		switch coin {
		case "BTC":
			priceVal, _, _, priceErr := jsonparser.Get(readBody, "USD")
			if priceErr == nil {
				log.Print("[++] BTC Price: " + string(priceVal))
			}
			retVal = CoinPrice{coin, "1", string(priceVal)}
		default:
			priceVal, _, _, eusd := jsonparser.Get(readBody, "USD")
			priceValBTC, _, _, ebtc := jsonparser.Get(readBody, "BTC")
			if eusd == nil && ebtc == nil {
				log.Printf("[++] %s Price: %s USD, %s BTC", coin, priceVal, priceValBTC)
			}
			retVal = CoinPrice{coin, string(priceValBTC), string(priceVal)}
		}
	} else {
		log.Print(err)
		retVal = CoinPrice{coin, "n/a", "n/a"}
	}

	return retVal
}
