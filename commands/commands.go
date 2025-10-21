package commands

import (
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/RhykerWells/asbwig/commands/economy"
	"github.com/RhykerWells/asbwig/commands/invite"
	"github.com/RhykerWells/asbwig/commands/moderation"
	"github.com/RhykerWells/asbwig/commands/ping"

	"github.com/RhykerWells/asbwig/commands/botOwner/banServer"
	"github.com/RhykerWells/asbwig/commands/botOwner/createInvite"
	"github.com/RhykerWells/asbwig/commands/botOwner/leaveServer"
	"github.com/RhykerWells/asbwig/commands/botOwner/setstatus"
	"github.com/RhykerWells/asbwig/commands/botOwner/unbanServer"
)

// InitCommands initializes the command handler, registers all
// available commands, and attaches the handler to the Discord session.
//
// After registration, the handler is connected to the session so that
// incoming message events are processed and routed to the correct
// command.
func InitCommands(session *discordgo.Session) {
	cmdHandler := dcommand.NewCommandHandler()

	cmdHandler.RegisterCommands(
		helpCmd,
		prefixCmd,

		ping.Command,
		invite.Command,

		banserver.Command,
		unbanserver.Command,
		createinvite.Command,
		leaveserver.Command,
		setstatus.Command,
	)

	economy.EconomySetup(cmdHandler)
	moderation.ModerationSetup(cmdHandler)
	session.AddHandler(cmdHandler.HandleMessageCreate)
}
