package main

import (
	"jimbot-go/lib"
	"time"
)

func main() {
	go lib.StartBot()
	for {
		time.Sleep(60 * time.Second)
	}
}
