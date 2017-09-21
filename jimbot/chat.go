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

// DecisionMaker : decide if a a reply is needed
func DecisionMaker() bool {
	timeSeed := time.Now().UnixNano()
	randNum := rand.Intn(int(timeSeed))
	log.Print("[***] RANDNUM = ", randNum)
	if randNum%10 == 0 {
		log.Println("[***] DECIDE TO RESPOND")
		return true
	}
	return false
}

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
	if strings.Contains(message, "谢谢") ||
		strings.Contains(message, "thanks") ||
		strings.Contains(message, "thank you") {
		return (HII + " np")
	}
	return message
}
