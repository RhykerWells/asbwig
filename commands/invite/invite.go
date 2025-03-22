package invite

import (
	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/Ranger-4297/asbwig/common"
	"github.com/Ranger-4297/asbwig/common/dcommand"
)

var Command = &dcommand.AsbwigCommand {
	Command:		[]string{"invite"},
	Description: 	"Creates an invite link for the bot",
	Run: (func(data *dcommand.Data) {
		functions.SendBasicMessage(data.Message.ChannelID, "[Invite link](<https://discord.com/oauth2/authorize?client_id=" + common.ConfigBotClientID + ">)")
	}),
}