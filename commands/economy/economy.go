package economy

const GuildEconomySchema = `
	CREATE TABLE IF NOT EXISTS economy_config (
		guild_id BIGINT PRIMARY KEY,
		max_bet BIGINT NOT NULL DEFAULT '5000',
		symbol TEXT NOT NULL DEFAULT '£',
		start_balance BIGINT NOT NULL DEFAULT '200'
	)
`