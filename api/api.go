package api

import (
	"github.com/Not-Cyrus/VoiceBot/commands"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	ID int
	DS *discordgo.Session
	BU *discordgo.User
}

func (b *Bot) Setup(token string) {
	b.DS, _ = discordgo.New(token)
	user, _ := b.DS.User("@me")
	b.BU = user
	b.DS.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0"
	route := commands.New()
	b.DS.AddHandler(route.MessageCreate)
	route.Add("summon", route.Summon)
	route.Add("playfile", route.Playfile)
	route.Add("play", route.PlayLink)
	route.Add("loop", route.Loop)
	//route.Add("resume", route.Resume)
	//route.Add("pause", route.Pause)
	route.Add("skip", route.Skip)
	route.Add("leave", route.Leave)
}

func (b *Bot) Run() {
	b.DS.Open()
}

func (b *Bot) Stop() {
	b.DS.Close()
}
