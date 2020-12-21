package commands

import (
	"strconv"

	"github.com/Not-Cyrus/VoiceBot/audio"
	"github.com/Not-Cyrus/VoiceBot/youtube"
	"github.com/bwmarrin/discordgo"
)

var S = NewSound() // really hacky code but LOL

type Sound struct {
	VC       *discordgo.VoiceConnection
	BChannel chan bool
	Playing  bool
}

func NewSound() *Sound {
	sound := &Sound{}
	return sound
}

func Getchannel(ds *discordgo.Session, guildID, authorID string) string {
	guild, _ := ds.State.Guild(guildID)
	for _, vs := range guild.VoiceStates {
		if vs.UserID == authorID {
			return vs.ChannelID
		}
	}
	return "No"
}

func (cmd *Commands) Playfile(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if len(ctx.Fields) > 0 {
		channel := Getchannel(s, m.GuildID, m.Author.ID)
		if channel == "No" {
			s.ChannelMessageSend(m.ChannelID, "Join a VC channel")
			return
		}
		var err error
		S.VC, err = s.ChannelVoiceJoin(m.GuildID, channel, false, false)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error joining your voice channel: "+err.Error())
			return
		}
		go S.Play("Assets/" + ctx.Fields[0] + ".mp3")
	}
}

func (cmd *Commands) PlayLink(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if len(ctx.Fields) > 0 {
		channel := Getchannel(s, m.GuildID, m.Author.ID)
		if channel == "No" {
			s.ChannelMessageSend(m.ChannelID, "Join a VC channel")
			return
		}
		var err error
		S.VC, err = s.ChannelVoiceJoin(m.GuildID, channel, false, false)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error joining your voice channel: "+err.Error())
			return
		}
		videoID, youtubeData, err := youtube.Search(ctx.Fields[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error getting YouTube data: "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Playing: "+videoID)
		go S.Play(string(youtubeData))
	}
}

func (S *Sound) Play(song string) {
	S.BChannel = make(chan bool, 1)
	S.Playing = true
	for {
		if S.VC != nil && S.Playing {
			audio.PlayAudioFile(S.VC, song, S.BChannel)
		} else {
			break
		}
	}
}

func (cmd *Commands) Skip(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	audio.Skip()
	S.Playing = false
}

func (cmd *Commands) Leave(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if S.VC != nil {
		S.BChannel <- true
		S.Playing = false
		S.VC.Disconnect()
		S.VC.Close()
	}
}
