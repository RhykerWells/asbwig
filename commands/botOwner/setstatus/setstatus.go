package setstatus

import (
	"strings"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
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
