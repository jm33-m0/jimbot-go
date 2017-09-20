package main

import (
	"jimbot-go/jimbot"
	"time"
)

func main() {
	go jimbot.StartBot()
	for {
		time.Sleep(60 * time.Second)
	}
}
