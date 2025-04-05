package jimbot

import (
	"log"
	"math/rand"

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
	randNum := rand.Int()
	if randNum%6 == 0 {
		log.Println("[+++] DECIDED TO RESPOND")
		return true
	}
	log.Println("[***] DECIDED TO IGNORE")
	return false
}

// ChoiceMaker : Select a random item from a slice
func ChoiceMaker(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
	return turing.GetResponse(message, InitConfig.OllamaModelName)
}
