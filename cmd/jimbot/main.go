package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jm33-m0/jimbot-go/jimbot"
)

func main() {
	go jimbot.StartBot()

	// Setup signal handling to exit on ctrl-c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-sigChan
	log.Printf("Received signal: %v, shutting down gracefully...", sig)
	os.Exit(0)
}
