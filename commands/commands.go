package commands

import (
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/Ranger-4297/asbwig/commands/help"
	"github.com/Ranger-4297/asbwig/commands/ping"
)

func InitCommands(session *discordgo.Session) {
	cmdHandler := dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		ping.Command,
		help.Command,
	)

	session.AddHandler(cmdHandler.HandleMessageCreate)
}