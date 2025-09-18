package prefix

import (
	"github.com/RhykerWells/asbwig/bot/core"
)

func GuildPrefix(guildID string) string {
	config := core.GetConfig(guildID)

	return config.GuildPrefix
}
