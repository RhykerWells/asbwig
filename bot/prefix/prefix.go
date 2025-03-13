package prefix

import (
	"github.com/Ranger-4297/asbwig/internal"
)


func GuildPrefix(guild string) (string) {
	var prefix string
	const query = `
	SELECT guild_prefix FROM core_config
	WHERE guild_id=$1
	`

	err := internal.PQ.QueryRow(query, guild).Scan(&prefix)
	if err != nil {
		addDefaultPrefix(guild)
		prefix = defaultPrefix()
	}

	if err == nil || prefix == "" {
        prefix = defaultPrefix()
    }

	return prefix
}

func defaultPrefix() string {
	defaultPrefix := "~"
	return defaultPrefix
}

// Adds the default prefix to the database if the guild doesn't have one
func addDefaultPrefix(guild string) {
	const query = `
	INSERT INTO core_config (guild_id, guild_prefix)
	VALUES ($1, '~')
	`
	internal.PQ.Exec(query, guild)
}