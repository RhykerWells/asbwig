package prefix

import (
	"github.com/RhykerWells/asbwig/common"
)

func GuildPrefix(guild string) string {
	var prefix string
	const query = `SELECT guild_prefix FROM core_config WHERE guild_id=$1`

	err := common.PQ.QueryRow(query, guild).Scan(&prefix)
	if err != nil {
		addDefaultPrefix(guild)
		prefix = "~"
	}

	return prefix
}

// Adds the default prefix to the database if the guild doesn't have one
func addDefaultPrefix(guild string) {
	const query = `INSERT INTO core_config (guild_id, guild_prefix) VALUES ($1, '~')`
	common.PQ.Exec(query, guild)
}

func ChangeGuildPrefix(guildID, prefix string) {
	const query = `UPDATE core_config SET guild_prefix = $1 WHERE guild_id = $2`
	common.PQ.Exec(query, prefix, guildID)
}