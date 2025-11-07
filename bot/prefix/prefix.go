package prefix

import (
	"github.com/RhykerWells/summit/bot/core"
)

// GuildPrefix returns the bots prefix for the current guild
func GuildPrefix(guildID string) string {
	config := core.GetConfig(guildID)

	return config.GuildPrefix
}
