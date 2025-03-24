package commands

import (
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/RhykerWells/asbwig/commands/invite"
	"github.com/RhykerWells/asbwig/commands/ping"
	"github.com/RhykerWells/asbwig/commands/economy"

	"github.com/RhykerWells/asbwig/commands/botOwner/eval"
	"github.com/RhykerWells/asbwig/commands/botOwner/setstatus"
)

var cmdHandler *dcommand.CommandHandler

func InitCommands(session *discordgo.Session) {
	cmdHandler = dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		helpCmd,

		ping.Command,
		invite.Command,

		setstatus.Command,
		eval.Command,
	)

	economySetup()
	session.AddHandler(cmdHandler.HandleMessageCreate)
} 

func economySetup() {
	common.InitSchema("Economy", economy.GuildEconomySchema)
}
