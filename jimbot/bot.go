package jimbot

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jm33-m0/jimbot-go/turing"
)

// bot api
var bot *tgbotapi.BotAPI

// userid of this bot
var botID int

// chat parameters
type chatParams struct {
	messageID     int
	chatID        int64
	userID        int64
	chatIsPrivate bool
	msgText       string
}

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

// InitConfig : cache config data
var InitConfig Config

// StartBot : Connect to Telegram bot API and start working
func StartBot() {
	// Save config in memory
	InitConfig = ReadConfig()

	// Login our bot
	loginToAPI()

	// Get updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("[-] Failed to get updates from Telegram server")
	}

	botID = bot.Self.ID
	for update := range updates {
		// handles empty update, prevent panic
		if update.Message == nil {
			continue
		}

		// handles each message
		go onMessage(update)
	}
}

func onMessage(update tgbotapi.Update) {
	/* on each message, do */
	// chat parameters
	var chat chatParams

	chat.chatID = update.Message.Chat.ID
	chat.messageID = update.Message.MessageID
	chat.chatIsPrivate = tgbotapi.Chat.IsPrivate(*update.Message.Chat)
	chat.msgText = update.Message.Text
	chat.userID = int64(update.Message.From.ID)

	log.Print("[**] Got msg from userID: ", chat.userID)
	// for strangers
	if chat.userID != InitConfig.BFID && chat.userID != InitConfig.GFID {
		log.Print("[!] Comparing userID <> BFID: ",
			chat.userID,
			" <> ",
			InitConfig.BFID,
			"\nStranger detected")
		onStranger(update, chat)
		return
	}

	if update.Message.IsCommand() {
		onCommand(update, chat)
	}

	// Write to histfile
	if WriteStringToFile("history.txt", "[*] "+chat.msgText, false) == nil {
		log.Println("[+] Message recorded")
	}

	// Mem dates
	memDate, greeting := checkMemDates()
	if _, err := os.Stat(".memdate_detected"); os.IsNotExist(err) {
		targetUserID := InitConfig.BFID
		if _, err = os.Stat(".mem4bf"); os.IsNotExist(err) {
			targetUserID = InitConfig.GFID
		}
		if memDate && chat.userID == targetUserID {

			// send photo with greeting
			_, err = bot.Send(tgbotapi.NewChatAction(chat.chatID, tgbotapi.ChatUploadPhoto))
			if err != nil {
				log.Println(err)
			}

			pic := tgbotapi.NewPhotoUpload(chat.chatID, "./img/mem.jpg")
			pic.Caption = greeting
			pic.ReplyToMessageID = chat.messageID
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
	} else if !memDate {
		if os.Remove(".memdate_detected") == nil {
			log.Print("[MEMDATE] not mem date, removing file")
		} else {
			log.Print("[MEM] Err deleting file")
		}
	}

	// from here, we handle normal messages
	var replyMsg tgbotapi.MessageConfig

	// get a reply
	replyProcessed := ProcessMsg(chat.msgText, chat.userID)
	// stop when we dont have a reply
	if replyProcessed == "" {
		return
	}

	// be quiet in group chats
	if chat.chatIsPrivate {
		replyMsg = tgbotapi.NewMessage(chat.chatID, replyProcessed)
	} else {

		// decide if make reponse
		if !DecisionMaker() {
			log.Println("[***] IGNORING MSG")
			return
		}

		log.Println("[***] MAKING RESPONSE")

		// Generate reply
		replyMsg = tgbotapi.NewMessage(chat.chatID, ProcessMsg(chat.msgText, chat.userID))

		// if not in private chat, quote msg
		replyMsg.ReplyToMessageID = chat.messageID
	}

	// send our reply
	_, err := bot.Send(tgbotapi.NewChatAction(chat.chatID, tgbotapi.ChatTyping))
	if err != nil {
		log.Println(err)
	}
	_, err = bot.Send(replyMsg)
	if err != nil {
		log.Println(err)
	}
}

// restrict access for strngers
func onStranger(update tgbotapi.Update, chat chatParams) {
	if chat.chatIsPrivate || isMentioned(update.Message) {
		// use turing 123
		turingReply := tgbotapi.NewMessage(chat.chatID, turing.GetResponse(chat.msgText))
		turingReply.ReplyToMessageID = chat.messageID
		_, err := bot.Send(tgbotapi.NewChatAction(chat.chatID, tgbotapi.ChatTyping))
		if err != nil {
			log.Println(err)
		}

		_, err = bot.Send(turingReply)
		if err != nil {
			log.Println(err)
		}
	}
	if update.Message.IsCommand() {
		onCommand(update, chat)
	}
}

func onCommand(update tgbotapi.Update, chat chatParams) {
	// bot commands
	cmd := update.Message.Command()
	cmd = strings.ToLower(cmd) // avoid markdown parsing in URL
	cmdMsg := tgbotapi.NewMessage(chat.chatID, "")
	cmdMsg.ReplyToMessageID = chat.messageID
	cmdArgs := update.Message.CommandArguments()

	// private commands
	privateCmds := [...]string{"greeting4mem", "pic4mem", "memdate", "count", "start"}
	if chat.userID != InitConfig.GFID && chat.userID != InitConfig.BFID {
		for _, priCmd := range privateCmds {
			log.Printf("[*] Private command: %s vs %s", cmd, priCmd)
			if cmd == priCmd {
				log.Printf("[!] Private command hit: %s", cmd)
				return
			}
		}
	}

	if cmd == "translate" {
		if update.Message.ReplyToMessage != nil {
			msgOrig := *update.Message.ReplyToMessage
			text := msgOrig.Text
			cmdMsg.ReplyToMessageID = msgOrig.MessageID
			cmdArgs = text
		} else {
			cmdArgs = strings.Join(
				strings.Split(
					update.Message.Text, " ")[1:], " ")
		}
	}

	if !strings.Contains(cmd, "google") &&
		!strings.Contains(cmd, "pic") {
		cmdMsg.ParseMode = "markdown"
	}

	cmdMsg.Text = ProcessCmd(cmd, cmdArgs, chat.userID)

	_, err := bot.Send(tgbotapi.NewChatAction(chat.chatID, tgbotapi.ChatTyping))
	if err != nil {
		log.Println(err)
	}

	_, err = bot.Send(cmdMsg)
	if err != nil {
		log.Println(err)
	}
}

// check if a message mentions the bot
func isMentioned(message *tgbotapi.Message) bool {
	reply2msg := message.ReplyToMessage
	if reply2msg == nil {
		return false
	}

	user := reply2msg.From.ID
	log.Printf("[+] reply2msg from: %d vs %d\n", user, botID)
	return user == botID
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
			retVal.GFID, _ = strconv.ParseInt(strings.Trim(value, "\n"), 0, 64)
		case "Boyfriend":
			retVal.BFName = value
		case "BFID":
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
	//	retVal)
	return retVal
}

func loginToAPI() {
	log.Print(InitConfig.Token)
	var err error
	bot, err = tgbotapi.NewBotAPI(InitConfig.Token)
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
