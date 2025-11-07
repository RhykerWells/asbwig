package invite

import (
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
)

var Command = &dcommand.SummitCommand{
	Command:     "invite",
	Category:    dcommand.CategoryGeneral,
	Description: "Creates an invite link for the bot",
	Run: (func(data *dcommand.Data) {
		functions.SendBasicMessage(data.ChannelID, "[Invite link](<https://discord.com/oauth2/authorize?client_id="+common.ConfigBotClientID+">)")
	}),
}
