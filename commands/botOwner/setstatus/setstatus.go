package setstatus

import (
	"strings"

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
	Args: []*dcommand.Args{
		{Name: "Status", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		functions.SetStatus(strings.Join(data.Args, " "))
		message := &discordgo.MessageSend{
			Content: "Status changed",
		}
		functions.SendMessage(data.ChannelID, message)
	}),
}
