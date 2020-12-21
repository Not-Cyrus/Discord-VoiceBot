package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Not-Cyrus/VoiceBot/api"
	"github.com/Not-Cyrus/VoiceBot/utils"
)

var (
	intid = 1
)

func loginbot(line string, id int) {
	bot := api.Bot{ID: id}
	bot.Setup(line)
	bot.Run()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bot.Stop()
}

func main() {
	tokens := utils.ReadLines("Tokens.txt")
	for index, line := range tokens {
		switch {
		case index == len(tokens)-1:
			loginbot(line, intid)
		default:
			go loginbot(line, intid)
		}
		intid++
	}
}
