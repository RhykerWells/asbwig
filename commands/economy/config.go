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
	Startbalance int64
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
		Startbalance: c.Startbalance,
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
		Startbalance: m.Startbalance,
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

func guildMemberAddToEconomy(m *discordgo.GuildMemberAdd) {
	config := GetConfig(m.GuildID)
	userEntry := models.EconomyUser{
		GuildID: config.GuildID,
		UserID:  m.User.ID,
		Cash:    config.Startbalance,
		Bank:    0,
	}
	userEntry.Insert(context.Background(), common.PQ, boil.Infer())
}

func guildMemberRemoveFromEconomy(m *discordgo.GuildMemberRemove) {
	models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyUserInventories(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
}