package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	Commands []*command
	Prefix   string
}

type command struct {
	Name string
	Run  handler
}

type Context struct {
	Fields  []string
	Content string
}

type handler func(*discordgo.Session, *discordgo.Message, *Context)

func New() *Commands {
	c := &Commands{}
	return c
}

func (cmds *Commands) Add(name string, fnc handler) *command {
	cmd := command{}
	cmd.Name = "c" + name
	cmd.Run = fnc
	cmds.Commands = append(cmds.Commands, &cmd)
	return &cmd
}

func (cmds *Commands) Match(m string) (*command, []string) {
	content := strings.Fields(m)
	if len(content) == 0 {
		return nil, nil
	}
	var c *command
	var rank int
	var commandKey int

	for commandKey, commandName := range content {
		for _, commandValue := range cmds.Commands {
			if commandValue.Name == commandName {
				return commandValue, content[commandKey:]
			}
			if strings.HasPrefix(commandValue.Name, commandName) {
				if len(commandName) > rank {
					c = commandValue
					rank = len(commandName)
				}
			}
		}
	}
	return c, content[commandKey:]
}

func (cmds *Commands) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == "yourid" {
		ctx := &Context{
			Content: strings.TrimSpace(m.Content),
		}
		cmd, fds := cmds.Match(ctx.Content)
		if cmd != nil {
			ctx.Fields = fds[1:]
			cmd.Run(s, m.Message, ctx)
			return
		}
	}
}
