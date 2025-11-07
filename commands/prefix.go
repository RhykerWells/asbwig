package commands

import (
	"github.com/RhykerWells/summit/bot/functions"
	prfx "github.com/RhykerWells/summit/bot/prefix"
	"github.com/RhykerWells/summit/common/dcommand"
)

var prefixCmd = &dcommand.SummitCommand{
	Command:     "prefix",
	Category:    dcommand.CategoryGeneral,
	Description: "Views the bot prefix",
	Args: []*dcommand.Args{
		{Name: "Prefix", Type: dcommand.String},
	},
	Run: (func(data *dcommand.Data) {
		prefix := prfx.GuildPrefix(data.GuildID)
		functions.SendBasicMessage(data.ChannelID, "This servers prefix is `"+prefix+"`")
	}),
}
