package commands

import (
	"github.com/RhykerWells/Summit/common/dcommand"
	"github.com/bwmarrin/discordgo"

	"github.com/RhykerWells/Summit/commands/economy"
	"github.com/RhykerWells/Summit/commands/invite"
	"github.com/RhykerWells/Summit/commands/moderation"
	"github.com/RhykerWells/Summit/commands/ping"

	banserver "github.com/RhykerWells/Summit/commands/botOwner/banServer"
	createinvite "github.com/RhykerWells/Summit/commands/botOwner/createInvite"
	leaveserver "github.com/RhykerWells/Summit/commands/botOwner/leaveServer"
	"github.com/RhykerWells/Summit/commands/botOwner/setstatus"
	unbanserver "github.com/RhykerWells/Summit/commands/botOwner/unbanServer"
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
