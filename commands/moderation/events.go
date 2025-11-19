package moderation

import (
	"context"

	"github.com/RhykerWells/Summit/bot/events"
	"github.com/RhykerWells/Summit/commands/moderation/models"
	"github.com/RhykerWells/Summit/common"
	"github.com/bwmarrin/discordgo"
)

// initEvents registers all the required event handlers to run on websocket events
func initEvents() {
	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate){
		guildAddModerationConfig,
	})

	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete){
		guildDeleteModerationConfig,
	})

	events.RegisterChannelCreatefunctions([]func(c *discordgo.ChannelCreate){
		func(c *discordgo.ChannelCreate) {
			refreshMuteSettingsOnChannel(GetConfig(c.GuildID), c.Channel)
		},
	})

	events.RegisterChannelUpdatefunctions([]func(c *discordgo.ChannelUpdate){
		func(c *discordgo.ChannelUpdate) {
			refreshMuteSettingsOnChannel(GetConfig(c.GuildID), c.Channel)
		},
	})
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
