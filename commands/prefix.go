package commands

import (
	"github.com/RhykerWells/asbwig/bot/functions"
	prfx "github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var prefixCmd = &dcommand.AsbwigCommand{
	Command:     "prefix",
	Category: 	 dcommand.CategoryGeneral,
	Description: "Views the bot prefix",
	Args: []*dcommand.Args{
		{Name: "Prefix", Type: dcommand.String},
	},
	Run: (func(data *dcommand.Data) {
		prefix := prfx.GuildPrefix(data.GuildID)
		functions.SendBasicMessage(data.ChannelID, "This servers prefix is `" + prefix + "`")
	}),
}