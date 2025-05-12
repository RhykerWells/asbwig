package commands

import (
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/RhykerWells/asbwig/commands/economy"
	"github.com/RhykerWells/asbwig/commands/invite"
	"github.com/RhykerWells/asbwig/commands/ping"

	"github.com/RhykerWells/asbwig/commands/botOwner/createInvite"
	"github.com/RhykerWells/asbwig/commands/botOwner/setstatus"
)

var cmdHandler *dcommand.CommandHandler

func InitCommands(session *discordgo.Session) {
	cmdHandler = dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		helpCmd,
		prefixCmd,

		ping.Command,
		invite.Command,

		createinvite.Command,
		setstatus.Command,
	)

	economy.EconomySetup(cmdHandler)
	session.AddHandler(cmdHandler.HandleMessageCreate)
}