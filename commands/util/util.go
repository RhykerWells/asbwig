package util

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/RhykerWells/asbwig/common/models"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

// IsGuildBanned returns a boolean of whether the guild is banned or not
func IsGuildBanned(guildID string) bool {
	exists, err := models.BannedGuilds(qm.Where("guild_id = ?", guildID)).Exists(context.Background(), common.PQ)
	if err != nil {
		return false
	}

	return exists
}