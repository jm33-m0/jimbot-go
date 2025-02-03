package jimbot

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jm33-m0/jimbot-go/huobi"
)

const (
	info    = "Under development"
	unknown = "Unknown command"
)

// ProcessCmd : handles bot commands
func ProcessCmd(command string, commandArgs string, userID int64) string {
	command = strings.ToLower(command)
	switch command {
	case "start":
		msg := start(userID)
		return msg
	case "remindmeto":
		// NOTE not finished
		// msg := remindMeTo(commandArgs)
		return info
	case "memdate":
		date := commandArgs
		log.Print("MemDay", "MemDay: "+date+"T00:00:00Z\n")
		err := UpdateConfig("MemDay", "MemDay: "+date+"T00:00:00Z")
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(".mem4bf")
		if err != nil {
			log.Println(err)
		}

		// setting memday for whom?
		if userID == InitConfig.GFID {
			if _, err := os.Create(".mem4bf"); err == nil {
				log.Print("[MEMDATE] mem4bf file created")
			} else {
				log.Print("[MEM] Err creating file mem4bf")
			}
		}

		return "Mem date set to " + date
	case "greeting4mem":
		msg := "Done"

		greeting := commandArgs
		log.Printf("Updating config with greeting: %s", greeting)
		err := UpdateConfig("MemdayGreetings", "MemdayGreetings: "+greeting)
		if err != nil {
			log.Println(err)
		}

		return msg
	case "pic4mem":
		msg := Search(commandArgs, true)
		msgSlice := strings.Split(msg, ": ")
		url := msgSlice[len(msgSlice)-1]
		log.Printf("Downloading from %s", url)

		if err := DownloadFile("./img/mem.jpg", url); err != nil {
			log.Println(err)
		}

		return msg
	case "count":
		return countMsg()
	case "translate":
		return ToEnglish(commandArgs)
	case "3_day_forecast":
		return info
	case "weather":
		if userID == InitConfig.GFID {
			return NowWeather(InitConfig.HisCity)
		}
		return NowWeather(InitConfig.HerCity)
	case "prices":
		msg := prices()
		return msg
	case "google":
		if commandArgs == "" {
			return unknown
		}
		return Search(commandArgs, false)
	case "pic":
		if commandArgs == "" {
			return unknown
		}
		return Search(commandArgs, true)
	case "huobi_market":
		return huobi.GetPrices()
	default:
		return unknown
	}
}

func start(userID int64) string {
	var msg string
	switch userID {
	case InitConfig.GFID:
		msg = "Hi, I'm your Telegram bot,\n"
		msg += "hope I'll be loved,\n"
		msg += "if not, well... blame him,"
		msg += "and... I love you two\n"
	case InitConfig.BFID:
		msg = "Hi, I'm your Telegram bot, and...\n"
		msg += "I'll always be here with you,\n"
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
	msg += strings.Repeat("-", 25)
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

// func googleSearch(query string, image bool) string {
//	log.Print("[###] Google query is : ", query)
//	return Search(query, image)
// }

func countMsg() string {
	counter := 0
	histfile, err := os.Open("history.txt")
	defer func() {
		err = histfile.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		log.Print("Failed to read history", err)
	}
	scanner := bufio.NewScanner(bufio.NewReader(histfile))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[*]") {
			counter++
		}
	}
	counter += 152843
	then, _ := time.Parse(time.RFC3339, "2016-12-05T14:23:00Z")
	log.Print("[TIME] then = ", then.Format(time.RFC3339), "\n")
	duration := time.Since(then)
	days := strconv.Itoa(int(int(duration.Hours()) / 24))
	log.Print("[TIME] duration = ", days)
	log.Print("[HIST LENGTH] ", counter)
	return (HII + " I've received " + strconv.Itoa(counter) + " messages from you two, and you've been together for " + days + " days.")
}
