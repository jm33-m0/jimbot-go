package jimbot

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// bot api
	bot *tgbotapi.BotAPI

	// chat parameters
	messageID     int
	chatID        int64
	userID        int64
	chatIsPrivate bool
	msgText       string
)

// Config : Read config info from text file
type Config struct {
	GFID int64
	BFID int64

	Token           string
	GFName          string
	BFName          string
	CSE             string
	HerCity         string
	HisCity         string
	MemDay          string
	MemdayGreetings string
	Birthday        string
	HuobiAccessKey  string
	HuobiSecretKey  string
}

// StartBot : Connect to Telegram bot API and start working
func StartBot() {
	// Login our bot
	loginToAPI()

	// Get updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("[-] Failed to get updates from Telegram server")
	}

	for update := range updates {
		// handles empty update, prevent panic
		if update.Message == nil {
			continue
		}

		// chat parameters
		chatID = update.Message.Chat.ID
		messageID = update.Message.MessageID
		chatIsPrivate = tgbotapi.Chat.IsPrivate(*update.Message.Chat)
		msgText = update.Message.Text
		userID = int64(update.Message.From.ID)

		// handles each message
		log.Print("[**] Got msg from userID: ", userID)
		go onMessage(update)
	}
}

func onMessage(update tgbotapi.Update) {
	/* on each message, do */

	// say no to strangers
	if userID != ReadConfig().BFID && userID != ReadConfig().GFID {
		log.Print("[!] Comparing userID <> BFID: ",
			userID,
			" <> ",
			ReadConfig().BFID)
		warningText := HUH + " I'm sorry, but I won't talk to you"
		warning := tgbotapi.NewMessage(chatID, warningText)
		_, err := bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
		if err != nil {
			log.Println(err)
		}

		_, err = bot.Send(warning)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// bot commands
	if update.Message.IsCommand() {
		cmd := update.Message.Command()
		cmd = strings.ToLower(cmd) // avoid markdown parsing in URL
		cmdMsg := tgbotapi.NewMessage(chatID, "")
		cmdMsg.ReplyToMessageID = messageID
		cmdArgs := update.Message.CommandArguments()
		cmdMsg.Text = ProcessCmd(cmd, cmdArgs, userID)
		if !strings.Contains(cmd, "google") &&
			!strings.Contains(cmd, "pic") {
			cmdMsg.ParseMode = "markdown"
		}
		_, err := bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
		if err != nil {
			log.Println(err)
		}

		_, err = bot.Send(cmdMsg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Write to histfile
	if WriteStringToFile("history.txt", "[*] "+msgText, false) == nil {
		log.Println("[+] Message recorded")
	}

	// Mem dates
	memDate, greeting := checkMemDates()
	if _, err := os.Stat(".memdate_detected"); os.IsNotExist(err) {
		targetUserID := ReadConfig().BFID
		if _, err = os.Stat(".mem4bf"); os.IsNotExist(err) {
			targetUserID = ReadConfig().GFID
		}
		if memDate && userID == targetUserID {

			// send photo with greeting
			_, err = bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatUploadPhoto))
			if err != nil {
				log.Println(err)
			}

			pic := tgbotapi.NewPhotoUpload(chatID, "./img/mem.jpg")
			pic.Caption = greeting
			pic.ReplyToMessageID = messageID
			_, err = bot.Send(pic)
			if err != nil {
				log.Println(err)
			}

			// mark done
			if _, err := os.Create(".memdate_detected"); err == nil {
				log.Print("[MEMDATE] MEM DAY! file created")
			} else {
				log.Print("[MEM] Err creating file")
			}
			return
		}
		log.Print("[MEM] No gf detected")
	} else if !memDate {
		if os.Remove(".memdate_detected") == nil {
			log.Print("[MEMDATE] not mem date, removing file")
		} else {
			log.Print("[MEM] Err deleting file")
		}
	}

	var replyMsg tgbotapi.MessageConfig
	if chatIsPrivate {
		replyMsg = tgbotapi.NewMessage(chatID, ProcessMsg(msgText, userID))
	} else {

		// decide if make reponse
		if !DecisionMaker() {
			log.Println("[***] IGNORING MSG")
			return
		}

		log.Println("[***] MAKING RESPONSE")

		// Generate reply
		replyMsg = tgbotapi.NewMessage(chatID, ProcessMsg(msgText, userID))

		// if not in private chat, quote msg
		replyMsg.ReplyToMessageID = messageID
	}

	// send our reply
	_, err := bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
	if err != nil {
		log.Println(err)
	}
	_, err = bot.Send(replyMsg)
	if err != nil {
		log.Println(err)
	}
}

// ReadConfig : Read config from config file
func ReadConfig() Config {
	var retVal Config
	lines, err := FileToLines("config.txt")
	if err != nil {
		log.Println("[-] Can't read config file")
		log.Fatal(err)
	}
	for _, line := range lines {
		value := strings.Split(line, ": ")[1]
		switch strings.Split(line, ": ")[0] {
		case "Girlfriend":
			retVal.GFName = value
		case "GFID":
			log.Print("[++] GFID string: ", value)
			retVal.GFID, _ = strconv.ParseInt(strings.Trim(value, "\n"), 0, 64)
		case "Boyfriend":
			retVal.BFName = value
		case "BFID":
			log.Print("[++] BFID string: ", value)
			retVal.BFID, _ = strconv.ParseInt(strings.Trim(value, "\n"), 0, 64)
		case "Token":
			retVal.Token = strings.Trim(value, "\n")
		case "CSE":
			retVal.CSE = strings.Trim(value, "\n")
		case "HerCity":
			retVal.HerCity = strings.Trim(value, "\n")
		case "HisCity":
			retVal.HisCity = strings.Trim(value, "\n")
		case "Birthday":
			retVal.Birthday = strings.Trim(value, "\n")
		case "MemDay":
			retVal.MemDay = strings.Trim(value, "\n")
		case "MemdayGreetings":
			retVal.MemdayGreetings = strings.Trim(value, "\n")
		case "HuobiAccessKey":
			retVal.HuobiAccessKey = strings.Trim(value, "\n")
		case "HuobiSecretKey":
			retVal.HuobiSecretKey = strings.Trim(value, "\n")
		default:
			log.Println("[-] Check your config file")
			os.Exit(1)
		}
	}
	// log.Print("======================Please check your config:======================\n",
	// 	retVal)
	return retVal
}

func loginToAPI() {
	log.Print(ReadConfig().Token)
	var err error
	bot, err = tgbotapi.NewBotAPI(ReadConfig().Token)
	if err != nil {
		log.Println("[-] Login failed, please check your token")
		log.Panic(err)
	}

	bot.Debug = true // for debugging

	log.Printf("[+] Authorized on account %s\n\n", bot.Self.UserName)
}

func checkMemDates() (bool, string) {
	birthDate, _ := time.Parse(time.RFC3339, ReadConfig().Birthday)
	anniversary, _ := time.Parse(time.RFC3339, ReadConfig().MemDay)
	nowDate := time.Now().Day()
	nowMonth := time.Now().Month()
	if (nowDate == birthDate.Day() &&
		nowMonth == birthDate.Month()) ||
		(nowDate == anniversary.Day() &&
			nowMonth == anniversary.Month()) {
		return true, ReadConfig().MemdayGreetings
	}
	return false, ""
}
