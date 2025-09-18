package core

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/core/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	// General
	GuildID     string
	GuildPrefix string
}

func (c *Config) ConfigToSQLModel() *models.CoreConfig {
	return &models.CoreConfig{
		GuildID:     c.GuildID,
		GuildPrefix: c.GuildPrefix,
	}
}

func ConfigFromModel(m *models.CoreConfig) *Config {
	return &Config{
		GuildID:     m.GuildID,
		GuildPrefix: m.GuildPrefix,
	}
}

func GetConfig(guildID string) *Config {
	model, err := models.FindCoreConfigG(context.Background(), guildID)
	if err == nil {
		return ConfigFromModel(model)
	}

	return &Config{
		GuildID:     guildID,
		GuildPrefix: "~",
	}
}

func SaveConfig(config *Config) error {
	err := config.ConfigToSQLModel().UpsertG(context.Background(), true, []string{"guild_id"}, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func DeleteConfig(config *Config) error {
	_, err := config.ConfigToSQLModel().Delete(context.Background(), common.PQ)
	if err != nil {
		return err
	}

	return nil
}

func guildAddCoreConfig(g *discordgo.GuildCreate) {
	config := GetConfig(g.ID)
	SaveConfig(config)
}

func guildDeleteCoreConfig(g *discordgo.GuildDelete) {
	config := GetConfig(g.ID)
	DeleteConfig(config)
}
