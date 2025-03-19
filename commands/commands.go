package commands

import (
	"github.com/bwmarrin/discordgo"
)

func InitCommands(session *discordgo.Session) {
	cmdHandler := newCommandHandler()

	cmdHandler.registerCommand(*Ping)
	cmdHandler.registerCommand(*Help)
	session.AddHandler(cmdHandler.handleMessage)
}