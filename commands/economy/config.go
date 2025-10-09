package economy

import (
	"context"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/aarondl/sqlboiler/v4/boil"
)

// Config defines the general struct to pass data to and from the dashboard template/context data
type Config struct {
	GuildID              string
	Min                  int64
	Max                  int64
	Maxbet               int64
	Symbol               string
	Startbalance         int64
	Customworkresponses  bool
	Customcrimeresponses bool
}

// ConfigToSQLModel converts a Config struct to the relevant SQLBoiler model
func (c *Config) ConfigToSQLModel() *models.EconomyConfig {
	return &models.EconomyConfig{
		GuildID:              c.GuildID,
		Min:                  c.Min,
		Max:                  c.Max,
		Maxbet:               c.Maxbet,
		Symbol:               c.Symbol,
		Startbalance:         c.Startbalance,
		Customworkresponses:  c.Customworkresponses,
		Customcrimeresponses: c.Customcrimeresponses,
	}
}

// ConfigFromModel converts the guild config SQLBoiler model to a Config struct 
func ConfigFromModel(m *models.EconomyConfig) *Config {
	return &Config{
		GuildID:              m.GuildID,
		Min:                  m.Min,
		Max:                  m.Max,
		Maxbet:               m.Maxbet,
		Symbol:               m.Symbol,
		Startbalance:         m.Startbalance,
		Customworkresponses:  m.Customworkresponses,
		Customcrimeresponses: m.Customcrimeresponses,
	}
}

// GetConfig returns the current or default guild config as a Config struct
func GetConfig(guildID string) *Config {
	model, err := models.FindEconomyConfigG(context.Background(), guildID)
	if err == nil {
		return ConfigFromModel(model)
	}

	return &Config{
		GuildID: guildID,
	}
}

// SaveConfig saves the passed Config struct via SQLBoiler
func SaveConfig(config *Config) error {
	err := config.ConfigToSQLModel().UpsertG(context.Background(), true, []string{models.EconomyConfigColumns.GuildID}, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
