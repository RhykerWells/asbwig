package commands

import (
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/RhykerWells/asbwig/commands/invite"
	"github.com/RhykerWells/asbwig/commands/ping"

	"github.com/RhykerWells/asbwig/commands/botOwner/eval"
	"github.com/RhykerWells/asbwig/commands/botOwner/setstatus"
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
