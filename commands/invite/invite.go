package invite

import (
	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "invite",
	Category: 	 dcommand.CategoryGeneral,
	Description: "Creates an invite link for the bot",
	Run: (func(data *dcommand.Data) {
		functions.SendBasicMessage(data.ChannelID, "[Invite link](<https://discord.com/oauth2/authorize?client_id="+common.ConfigBotClientID+">)")
	}),
}
