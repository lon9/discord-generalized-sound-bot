package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		panic(errors.New("CONFIG_FILE is not defined"))
	}
	config, err := NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	b, err := NewBot(config)
	if err != nil {
		panic(err)
	}

	defer b.Close()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Closing sessions.")
}
