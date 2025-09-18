package moderation

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
)

//go:generate sqlboiler --no-hooks psql

func ModerationSetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Moderation", GuildModerationSchema...)

	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate) {
		guildAddModerationConfig,
	})
	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete) {
		guildDeleteModerationConfig,
	})

	initWeb()

	scheduleAllPendingUnmutes()
	scheduleAllPendingUnbans()

	cmdHandler.RegisterCommands(
		warnCommand,
		muteCommand,
		unmuteCommand,
		kickCommand,
		banCommand,
		unbanCommand,
	)
}

func guildAddModerationConfig(g *discordgo.GuildCreate) {
	config := GetConfig(g.ID)
	SaveConfig(config)
}

func guildDeleteModerationConfig(g *discordgo.GuildDelete) {
	config, err := models.ModerationConfigs(qm.Where("guild_id = ?", g.ID)).One(context.Background(), common.PQ)
	if err != nil {
		return
	}

	config.Delete(context.Background(), common.PQ)
}
