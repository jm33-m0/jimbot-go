package jimbot

import (
	"log"
	"math/rand"
	"strings"
	"time"
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

// DecisionMaker : decide if a a reply is needed, randomly
func DecisionMaker() bool {
	timeSeed := time.Now().UnixNano()
	randNum := rand.Intn(int(timeSeed))
	log.Print("[***] RANDNUM = ", randNum)
	if randNum%5 == 0 {
		log.Println("[***] DECIDED TO RESPOND")
		return true
	}
	return false
}

// ChoiceMaker : Select a random item from a slice
func ChoiceMaker(choices []string) string {
	rand.Seed(time.Now().Unix())
	return choices[rand.Intn[len(choices)]]
}

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
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
		if DecisionMaker() {
			return "没错"
		}
		return "不存在的"
	}
	return HUH
}
