package main

import (
	"time"

	//"jimbot-go/jimbot"
	"github.com/jm33-m0/jimbot-go/jimbot"
)

func main() {
	go jimbot.StartBot()
	for {
		time.Sleep(60 * time.Second)
	}
}
