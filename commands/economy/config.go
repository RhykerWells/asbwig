package economy

import (
	"context"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	// General
	GuildID	string
	Min int64
	Max int64
	Maxbet int64
	Symbol string
	Customworkresponses bool
	Customcrimeresponses bool
}

func (c *Config) ConfigToSQLModel() *models.EconomyConfig {
	return &models.EconomyConfig{
		GuildID: c.GuildID,
		Min: c.Min,
		Max: c.Max,
		Maxbet: c.Maxbet,
		Symbol: c.Symbol,
		Customworkresponses: c.Customworkresponses,
		Customcrimeresponses: c.Customcrimeresponses,
	}
}

func ConfigFromModel(m *models.EconomyConfig) *Config {
	return &Config{
		GuildID: m.GuildID,
		Min: m.Min,
		Max: m.Max,
		Maxbet: m.Maxbet,
		Symbol: m.Symbol,
		Customworkresponses: m.Customworkresponses,
		Customcrimeresponses: m.Customcrimeresponses,
	}
}

func GetConfig(guildID string) *Config {
	model, err := models.FindEconomyConfigG(context.Background(), guildID)
	if err == nil {
		return ConfigFromModel(model)
	}

	return &Config{
		GuildID: guildID,
	}
}

func SaveConfig(config *Config) error {
	err := config.ConfigToSQLModel().UpsertG(context.Background(), true, []string{"guild_id"}, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func guildAddEconomyConfig(g *discordgo.GuildCreate) {
	config := GetConfig(g.ID)
	SaveConfig(config)
}

func guildDeleteEconomyConfig(g *discordgo.GuildDelete) {
	config, err := models.EconomyConfigs(qm.Where("guild_id = ?", g.ID)).One(context.Background(), common.PQ)
	if err != nil {
		return
	}

	config.Delete(context.Background(), common.PQ)
}