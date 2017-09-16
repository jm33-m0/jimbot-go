package lib

import (
	"jimbot-go/util"
	"strings"
)

// ProcessCmd : handles bot commands
func ProcessCmd(command string, userID int64) string {
	switch command {
	case "start":
		msg := start(userID)
		return msg
	case "google":
		return "Under development"
	case "stat":
		return "Under development"
	case "translate":
		return "Under development"
	case "pic":
		return "Under development"
	case "3_day_forecast":
		return "Under development"
	case "weather":
		return "Under development"
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
	btcPrice := util.GetPrice("BTC")
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
	coinPrice := util.GetPrice(coin)
	var msg string
	msg += coinPrice.CoinName + " -> USD : " + coinPrice.PriceInUSD + "\n"
	msg += coinPrice.CoinName + " -> BTC : " + coinPrice.PriceInBTC + "\n"
	return msg
}
