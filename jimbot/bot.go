package jimbot

import (
	"encoding/json" // new import
	"log"
	"os"
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

// Updated Config : now matches config.json keys with json tags
type Config struct {
	GFID            int64  `json:"GFID"`
	BFID            int64  `json:"BFID"`
	Token           string `json:"Token"`
	GFName          string `json:"Girlfriend"`
	BFName          string `json:"Boyfriend"`
	CSE             string `json:"CSE"`
	HerCity         string `json:"HerCity"`
	HisCity         string `json:"HisCity"`
	MemDay          string `json:"MemDay"`
	MemdayGreetings string `json:"MemdayGreetings"`
	Birthday        string `json:"Birthday"`
	HuobiAccessKey  string `json:"HuobiAccessKey"`
	HuobiSecretKey  string `json:"HuobiSecretKey"`
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
		chatbot(update, chat)
		return
	}

	// for BF and GF

	if update.Message.IsCommand() {
		onCommand(update, chat)
	} else {
		chatbot(update, chat)
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

// restrict access for strangers
func chatbot(update tgbotapi.Update, chat chatParams) {
	if chat.chatIsPrivate || isMentioned(update.Message) {
		resp := turing.GetResponse(update.Message.Text)
		msg := tgbotapi.NewMessage(chat.chatID, resp)
		msg.ReplyToMessageID = chat.messageID
		_, err := bot.Send(tgbotapi.NewChatAction(chat.chatID, tgbotapi.ChatTyping))
		if err != nil {
			log.Println(err)
		}
		_, err = bot.Send(msg)
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

// Modified ReadConfig: loads configuration from config.json file
func ReadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("[-] Can't read config file: ", err)
	}
	defer file.Close()

	var config Config
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal("[-] Failed to decode config: ", err)
	}
	return config
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
