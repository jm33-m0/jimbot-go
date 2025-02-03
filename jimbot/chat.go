package jimbot

import (
	"log"
	"math/rand"
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

// DecisionMaker : decide if a a reply is needed, randomly
func DecisionMaker() bool {
	timeSeed := time.Now().UnixNano()
	randNum := rand.Intn(int(timeSeed))
	if randNum%12 == 0 {
		log.Println("[***] DECIDED TO RESPOND")
		return true
	}
	return false
}

// ChoiceMaker : Select a random item from a slice
func ChoiceMaker(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
	// say something
	if DecisionMaker() {
		return turing.GetResponse(message, InitConfig.OllamaModelName)
	}
	return ""
}
