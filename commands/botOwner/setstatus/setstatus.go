package setstatus

import (
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/commands/util"
	"github.com/RhykerWells/summit/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.SummitCommand{
	Command:      "setstatus",
	Category:     dcommand.CategoryOwner,
	Description:  "Changes the bot status",
	ArgsRequired: 1,
	Args: []*dcommand.Arg{
		{Name: "Status", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		functions.SetStatus(data.ParsedArgs[0].String())
		message := &discordgo.MessageSend{
			Content: "Status changed",
		}
		functions.SendMessage(data.ChannelID, message)
	}),
}
