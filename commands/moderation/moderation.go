package moderation

import (
	"errors"

	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

//go:generate sqlboiler --no-hooks psql

func ModerationSetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Moderation", GuildModerationSchema...)
	initWeb()
	scheduleAllPendingUnmutes()
	cmdHandler.RegisterCommands(
		warnCommand,
		muteCommand,
		unmuteCommand,
		kickCommand,
		banCommand,
		unbanCommand,
	)
}

func GuildModerationAdd(guild_id string) {
	const query = `SELECT guild_id FROM moderation_config WHERE guild_id=$1`
	err := common.PQ.QueryRow(query, guild_id)
	if err != nil {
		guildModerationDefault(guild_id)
	}
}

func guildModerationDefault(guild_id string) {
	const query = `INSERT INTO moderation_config (guild_id) VALUES ($1)`
	common.PQ.Exec(query, guild_id)
}

func getGuildModLogChannel(guildID string) (string, error) {
	var logChannel string
	query := `SELECT mod_log FROM moderation_config WHERE guild_id=$1`

	err := common.PQ.QueryRow(query, guildID).Scan(&logChannel)
	if err != nil {
		return "", errors.New("no modlog channel found")
	}

	ok := util.HasPerms(guildID, logChannel, common.Bot.ID, discordgo.PermissionSendMessages)
	if !ok {
		return "", errors.New("cannot send message in the modlog channel")
	}

	return logChannel, nil
}