package jimbot

import (
	"strings"
)

const info = "Under development"

// ProcessCmd : handles bot commands
func ProcessCmd(command string, userID int64) string {
	switch command {
	case "start":
		msg := start(userID)
		return msg
	case "google":
		return info
	case "stat":
		return info
	case "translate":
		return info
	case "pic":
		return info
	case "3_day_forecast":
		return info
	case "weather":
		return info
	case "prices":
		msg := prices()
		return msg
	default:
		return "Unknown command"
	}
}

func start(userID int64) string {
	var msg string
	switch userID {
	case ReadConfig().GFID:
		msg = "Hi, I'm your Telegram bot,\n"
		msg += "hope you like me\n"
		msg += "if not, well... blame him"
		msg += "and... I love you too\n"
	case ReadConfig().BFID:
		msg = "Hi, I guess you already knew me well\n"
		msg += "let's cut the bullshit\n"
		msg += "and... I love you two\n"
	default:
		msg = "There must be something wrong...\n"
	}
	msg += KISS
	return msg
}

func prices() string {
	btcPrice := GetPrice("BTC")
	msg := HII + " I got this list\n`"
	msg += strings.Repeat("-", 35)
	msg += "\n"
	msg += btcPrice.CoinName + " -> USD : " + btcPrice.PriceInUSD + "\n"
	msg += getAltcoinPrices("XMR")
	msg += getAltcoinPrices("ETH")
	msg += getAltcoinPrices("ETC")
	msg += getAltcoinPrices("BCC")
	msg += "`"
	return msg
}

func getAltcoinPrices(coin string) string {
	coinPrice := GetPrice(coin)
	var msg string
	msg += coinPrice.CoinName + " -> USD : " + coinPrice.PriceInUSD + "\n"
	msg += coinPrice.CoinName + " -> BTC : " + coinPrice.PriceInBTC + "\n"
	return msg
}
