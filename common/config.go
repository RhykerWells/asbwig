package common

import (
	"context"
	"os"
	"strings"

	"github.com/RhykerWells/asbwig/common/models"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
)

var (
	ConfigBotName     = os.Getenv("ASBWIG_BOTNAME")
	ConfigBotToken    = os.Getenv("ASBWIG_TOKEN")
	ConfigBotClientID = os.Getenv("ASBWIG_CLIENTID")
	ConfigBotSecret   = os.Getenv("ASBWIG_CLIENTSECRET")
	ConfigASBWIGHost  = os.Getenv("ASBWIG_HOST")

	ConfigPGHost     = os.Getenv("ASBWIG_PGHOST")
	ConfigPGDB       = os.Getenv("ASBWIG_PGDB")
	ConfigPGUsername = os.Getenv("ASBWIG_PGUSER")
	ConfigPGPassword = os.Getenv("ASBWIG_PGPASSWORD")

	ConfigTermsURLOverride = os.Getenv("ASBWIG_TERMSURLOVERRIDE")
	ConfigPrivacyURLOverride = os.Getenv("ASBWIG_PRIVACYURLOVERRIDE")

	ConfigBotOwner = os.Getenv("ASBWIG_OWNERID")
)

func ConfigDgoBotToken() string {
	token := ConfigBotToken
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}
	return token
}

type Config struct {
	// General
	GuildID	string
	GuildPrefix string
}

func (c *Config) ConfigToSQLModel() *models.CoreConfig {
	return &models.CoreConfig{
		GuildID: c.GuildID,
		GuildPrefix: c.GuildPrefix,
	}
}

func ConfigFromModel(m *models.CoreConfig) *Config {
	return &Config{
		GuildID: m.GuildID,
		GuildPrefix: m.GuildPrefix,
	}
}

func GetConfig(guildID string) *Config {
	model, err := models.FindCoreConfigG(context.Background(), guildID)
	if err == nil {
		return ConfigFromModel(model)
	}

	return &Config{
		GuildID: guildID,
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
	_, err := config.ConfigToSQLModel().Delete(context.Background(), PQ)
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