package moderation

import (
	"context"

	"github.com/RhykerWells/summit/bot/events"
	"github.com/RhykerWells/summit/commands/moderation/models"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

//go:generate sqlboiler --no-hooks psql

// ModerationSetup runs the following:
//   - The schema initialiser
//   - Registration of the guild join/leave functions
//   - Initialises the web plugin
//   - Initialises any other required middlewares
//   - Registration of the moderation commands & their pagination
func ModerationSetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Moderation", GuildModerationSchema...)

	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate){
		guildAddModerationConfig,
	})
	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete){
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

// guildAddModerationConfig creates the intial configs for the moderation system for a specified guild
func guildAddModerationConfig(g *discordgo.GuildCreate) {
	config := GetConfig(g.ID)
	SaveConfig(config)
}

// guildDeleteModerationConfig deletes the config for the moderation system for a specified guild
func guildDeleteModerationConfig(g *discordgo.GuildDelete) {
	config, err := models.ModerationConfigs(models.ModerationConfigWhere.GuildID.EQ(g.ID)).One(context.Background(), common.PQ)
	if err != nil {
		return
	}

	config.Delete(context.Background(), common.PQ)
}
