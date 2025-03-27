package economy

var GuildEconomySchema = []string{`
CREATE TABLE IF NOT EXISTS economy_config (
	guild_id TEXT PRIMARY KEY,
	maxbet BIGINT NOT NULL DEFAULT '5000',
	symbol TEXT NOT NULL DEFAULT '£',
	startbalance BIGINT NOT NULL DEFAULT '200'
);
`,`
CREATE TABLE IF NOT EXISTS economy_cash (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	cash BIGINT NOT NULL
);
`,`
CREATE TABLE IF NOT EXISTS economy_bank (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	balance BIGINT NOT NULL
)`,
}