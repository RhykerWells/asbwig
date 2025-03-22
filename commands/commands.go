package commands

import (
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/Ranger-4297/asbwig/commands/ping"

	"github.com/Ranger-4297/asbwig/commands/botOwner/setstatus"
)

func InitCommands(session *discordgo.Session) {
	cmdHandler := dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		helpCmd,

		ping.Command,

		setstatus.Command,
	)

	session.AddHandler(cmdHandler.HandleMessageCreate)
}