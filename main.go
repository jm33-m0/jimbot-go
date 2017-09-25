package main

import (
	"jimbot-go/jimbot"
	"time"
	//"github.com/jm33-m0/jimbot-go/jimbot"
)

func main() {
	go jimbot.StartBot()
	for {
		time.Sleep(60 * time.Second)
	}
}
