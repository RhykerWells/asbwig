package economy

import (
	"context"

	"github.com/RhykerWells/Summit/commands/economy/models"
	"github.com/RhykerWells/Summit/common"
	"github.com/aarondl/sqlboiler/v4/boil"
)

// Config defines the general struct to pass data to and from the dashboard template/context data
type Config struct {
	// General
	GuildID             string
	EconomyEnabled      bool
	EconomySymbol       string
	EconomyStartBalance int64

	// Money making management
	EconomyMinReturn int64
	EconomyMaxReturn int64
	EconomyMaxBet    int64

	// Custom responses
	EconomyCustomWorkResponsesEnabled  bool
	EconomyCustomWorkResponses         []string
	EconomyCustomCrimeResponsesEnabled bool
	EconomyCustomCrimeResponses        []string
}

// ConfigToSQLModel converts a Config struct to the relevant SQLBoiler model
func (c *Config) ConfigToSQLModel() *models.EconomyConfig {
	return &models.EconomyConfig{
		// General
		GuildID:             c.GuildID,
		EconomyEnabled:      c.EconomyEnabled,
		EconomySymbol:       c.EconomySymbol,
		EconomyStartBalance: c.EconomyStartBalance,

		// Money making management
		EconomyMinReturn: c.EconomyMinReturn,
		EconomyMaxReturn: c.EconomyMaxReturn,
		EconomyMaxBet:    c.EconomyMaxBet,

		// Custom responses
		EconomyCustomWorkResponsesEnabled:  c.EconomyCustomWorkResponsesEnabled,
		EconomyCustomWorkResponses:         c.EconomyCustomWorkResponses,
		EconomyCustomCrimeResponsesEnabled: c.EconomyCustomCrimeResponsesEnabled,
		EconomyCustomCrimeResponses:        c.EconomyCustomCrimeResponses,
	}
}

// ConfigFromModel converts the guild config SQLBoiler model to a Config struct
func ConfigFromModel(m *models.EconomyConfig) *Config {
	return &Config{
		// General
		GuildID:             m.GuildID,
		EconomyEnabled:      m.EconomyEnabled,
		EconomySymbol:       m.EconomySymbol,
		EconomyStartBalance: m.EconomyStartBalance,

		// Money making management
		EconomyMinReturn: m.EconomyMinReturn,
		EconomyMaxReturn: m.EconomyMaxReturn,
		EconomyMaxBet:    m.EconomyMaxBet,

		// Custom responses
		EconomyCustomWorkResponsesEnabled:  m.EconomyCustomWorkResponsesEnabled,
		EconomyCustomWorkResponses:         m.EconomyCustomWorkResponses,
		EconomyCustomCrimeResponsesEnabled: m.EconomyCustomCrimeResponsesEnabled,
		EconomyCustomCrimeResponses:        m.EconomyCustomCrimeResponses,
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

// getGuildShop returns the guild shop
func getGuildShop(guildID string) models.EconomyShopSlice {
	models, err := models.EconomyShops(models.EconomyShopWhere.GuildID.EQ(guildID)).All(context.Background(), common.PQ)
	if err != nil {
		return nil
	}

	return models
}
