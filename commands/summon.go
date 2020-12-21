package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (cmd *Commands) Summon(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if len(ctx.Fields) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Put in an invite you stupid fuck")
		return
	}
	_, err := s.InviteAccept(ctx.Fields[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Joined")
}
