package util

import (
	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

func OwnerCommand(inner dcommand.Run) dcommand.Run {
	return func(data *dcommand.Data) {
		if data.Author.ID == common.ConfigBotOwner {
			inner(data)
		} else {
			functions.SendBasicMessage(data.ChannelID, "This is a bot-owner only command.")
		}
	}
}

func AdminOrManageServerCommand(inner dcommand.Run) dcommand.Run {
	return func(data *dcommand.Data) {
		perms, _ := data.Session.State.UserChannelPermissions(data.Author.ID, data.ChannelID)
		if perms&discordgo.PermissionAdministrator == 8 || perms&discordgo.PermissionManageServer == 32 {
			inner(data)
		} else {
			functions.SendBasicMessage(data.ChannelID, "You need `Administrator` or `ManageServer` permissions to use this command.")
		}
	}
}
