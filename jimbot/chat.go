package jimbot

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jm33-m0/jimbot-go/turing"
)

// Emojis
const (
	KISS     = "😘"
	HEART    = "💕"
	TONGUE   = "😋"
	UNHAPPY  = "😐"
	SILENT   = "😶"
	HUH      = "🌚"
	HII      = "🙃"
	SWEAT    = "😓"
	SURPRISE = "😮"
)

var (
	emojis = make([]string, 0)
	yesNo  = make([]string, 0)
	what   = make([]string, 0)
)

// DecisionMaker : decide if a a reply is needed, randomly
func DecisionMaker() bool {
	timeSeed := time.Now().UnixNano()
	randNum := rand.Intn(int(timeSeed))
	log.Print("[***] RANDNUM = ", randNum)
	if randNum%12 == 0 {
		log.Println("[***] DECIDED TO RESPOND")
		return true
	}
	return false
}

// ChoiceMaker : Select a random item from a slice
func ChoiceMaker(choices []string) string {
	rand.Seed(time.Now().Unix())
	return choices[rand.Intn(len(choices))]
}

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
	// emojis for response
	emojis = append(emojis,
		KISS,
		HEART,
		TONGUE,
		UNHAPPY,
		SILENT,
		HUH,
		SWEAT,
		HII,
		SURPRISE)

	// answers for yes or no
	yesNo = append(yesNo,
		"不存在的",
		"嗯哼",
		"说的没错",
		"不对不对不对～",
		"nope",
		"no way",
		"dunno",
		"yeha",
		"yea",
		"yeah",
		"ok")

	// answers for what
	what = append(what,
		"不知道",
		"dunno",
		"emmm",
		"what?",
		"ask google")

	message = strings.ToLower(message)

	if strings.Contains(message, "谢谢") ||
		strings.Contains(message, "thanks") ||
		strings.Contains(message, "thank you") {
		return (HII + " np")
	} else if strings.Contains(message, "good night") ||
		strings.Contains(message, "goodnight") ||
		strings.Contains(message, "晚安") {
		return (KISS + " Good night!")
	} else if strings.Contains(message, "jimbot") ||
		strings.Contains(message, "jim bot") {
		return (HII + " huh?")
	} else if strings.Contains(message, "是不是") ||
		strings.Contains(message, "是吗") ||
		strings.Contains(message, "是么") ||
		strings.Contains(message, "对不") ||
		strings.Contains(message, "对吗") ||
		strings.Contains(message, "对么") {
		return ChoiceMaker(yesNo)
	} else if strings.Contains(message, "是啥") ||
		strings.Contains(message, "是什么") ||
		strings.Contains(message, "什么") {
		if DecisionMaker() {
			return Search(message, DecisionMaker())
		}
		return ChoiceMaker(what)
	} else if strings.HasPrefix(message, "google") {
		q := strings.Split(message, "google")[1:]
		query := strings.Join(q, " ")
		return Search(query, false)
	}

	// say something
	if DecisionMaker() {
		return turing.GetResponse(message)
	}

	return ChoiceMaker(emojis)
}
