package commands

import (
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/Ranger-4297/asbwig/commands/invite"
	"github.com/Ranger-4297/asbwig/commands/ping"

	"github.com/Ranger-4297/asbwig/commands/botOwner/eval"
	"github.com/Ranger-4297/asbwig/commands/botOwner/setstatus"
)

func InitCommands(session *discordgo.Session) {
	cmdHandler := dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		helpCmd,

		ping.Command,
		invite.Command,

		setstatus.Command,
		eval.Command,
	)

	session.AddHandler(cmdHandler.HandleMessageCreate)
}