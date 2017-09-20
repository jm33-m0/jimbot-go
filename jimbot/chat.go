package jimbot

import (
	"strings"
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

// ProcessMsg : handles chat messages
func ProcessMsg(message string, userID int64) string {
	if strings.Contains(message, "谢谢") ||
		strings.Contains(message, "thanks") ||
		strings.Contains(message, "thank you") {
		return (HII + " np")
	}
	return message
}
