package economy

var GuildEconomySchema = []string{`
CREATE TABLE IF NOT EXISTS economy_config (
	guild_id TEXT PRIMARY KEY,
	min BIGINT NOT NULL DEFAULT '200',
	MAX BIGINT NOT NULL DEFAULT '500',
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
);`,`
CREATE TABLE IF NOT EXISTS economy_cooldowns (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	type TEXT NOT NULL,
	expires_at TIMESTAMP
);`,`
DO $$
BEGIN
  BEGIN
    ALTER TABLE economy_cooldowns ADD CONSTRAINT economy_cooldowns_unique UNIQUE (guild_id, user_id, type);
  EXCEPTION
    WHEN duplicate_table THEN RAISE NOTICE 'Table constraint economy_cooldowns_unique already exists';
  END;
END $$;
`,
}