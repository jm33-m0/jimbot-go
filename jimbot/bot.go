package jimbot

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

// Config : Read config info from text file
type Config struct {
	Token  string
	GFName string
	BFName string
	CSE    string
	GFID   int64
	BFID   int64
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
		// chat parameters
		chatID := update.Message.Chat.ID
		messageID := update.Message.MessageID
		msgText := update.Message.Text
		chatIsPrivate := tgbotapi.Chat.IsPrivate(*update.Message.Chat)
		userID := int64(update.Message.From.ID)
		log.Print("[**] Got msg from userID: ", userID)

		if update.Message == nil {
			continue
		} else if userID != ReadConfig().BFID && userID != ReadConfig().GFID {
			log.Print("[!] Comparing userID <> BFID: ",
				userID,
				" <> ",
				ReadConfig().BFID)
			warningText := HUH + " I'm sorry, but I won't talk to you"
			warning := tgbotapi.NewMessage(chatID, warningText)
			bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
			bot.Send(warning)
			continue
		} else if update.Message.IsCommand() {
			cmd := update.Message.Command()
			cmdMsg := tgbotapi.NewMessage(chatID, "")
			cmdMsg.ReplyToMessageID = messageID
			cmdMsg.ParseMode = "markdown"
			// cmdMsg.DisableWebPagePreview = true TODO : Get file directly from URL and upload it
			cmdArgs := update.Message.CommandArguments()
			cmdMsg.Text = ProcessCmd(cmd, cmdArgs, userID)
			bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
			bot.Send(cmdMsg)
			continue
		}

		// Write to histfile
		f, err := os.Create("history.txt")
		if err != nil {
			log.Print("[==] failed to create histfile")
		}
		defer f.Close()
		f.WriteString("[*] " + update.Message.Text)

		// decide if make reponse
		if !DecisionMaker() {
			log.Println("[***] IGNORING MSG")
			continue
		}

		log.Println("[***] MAKING RESPONSE")

		// Generate reply
		replyMsg := tgbotapi.NewMessage(chatID, ProcessMsg(msgText, userID))

		// if not in private chat, quote msg
		if !chatIsPrivate {
			replyMsg.ReplyToMessageID = messageID
		}

		bot.Send(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
		bot.Send(replyMsg)
	}
}

// FileToLines : Read lines from a text file
func FileToLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
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
		default:
			log.Println("[-] Check your config file")
			os.Exit(1)
		}
	}
	log.Print("======================Please check your config:======================\n",
		retVal)
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
