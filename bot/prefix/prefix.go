package prefix

import (
	"github.com/RhykerWells/asbwig/bot/core"
)

// GuildPrefix returns the bots prefix for the current guild
func GuildPrefix(guildID string) string {
	config := core.GetConfig(guildID)

	return config.GuildPrefix
}
