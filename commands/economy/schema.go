package economy

var GuildEconomySchema = []string{`
CREATE TABLE IF NOT EXISTS economy_config (
	guild_id TEXT PRIMARY KEY,
	min BIGINT NOT NULL DEFAULT 200,
	max BIGINT NOT NULL DEFAULT 500,
	maxbet BIGINT NOT NULL DEFAULT 5000,
	symbol TEXT NOT NULL DEFAULT '£',
	startbalance BIGINT NOT NULL DEFAULT 200,
	customworkresponses BOOLEAN NOT NULL DEFAULT 'false',
	customcrimeresponses BOOLEAN NOT NULL DEFAULT 'false'
);
`,`
CREATE TABLE IF NOT EXISTS economy_custom_responses (
    guild_id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
    response TEXT NOT NULL,
    CONSTRAINT fk_guild_work_responses FOREIGN KEY (guild_id)
        REFERENCES economy_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE TABLE IF NOT EXISTS economy_users (
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	cash BIGINT NOT NULL,
	bank BIGINT NOT NULL,
	PRIMARY KEY (guild_id, user_id),
	CONSTRAINT fk_guild_user FOREIGN KEY (guild_id)
        REFERENCES economy_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE INDEX IF NOT EXISTS idx_guild_users
	ON economy_users (guild_id, user_id);
`,`
CREATE TABLE IF NOT EXISTS economy_cooldowns (
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	type TEXT NOT NULL,
	expires_at TIMESTAMP,
	PRIMARY KEY (guild_id, user_id, type),
	CONSTRAINT fk_guild_user_cooldown FOREIGN KEY (guild_id)
        REFERENCES economy_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE INDEX IF NOT EXISTS idx_guild_user_cooldowns
	ON economy_cooldowns (guild_id, user_id);
`,`
CREATE TABLE IF NOT EXISTS economy_shop (
    guild_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    quantity BIGINT,
    role TEXT,
    reply TEXT NOT NULL,
    soldby TEXT,
    PRIMARY KEY (guild_id, name),
	CONSTRAINT fk_guild_shop FOREIGN KEY (guild_id)
		REFERENCES economy_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE TABLE IF NOT EXISTS economy_createitem (
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	name TEXT,
	description TEXT,
	price BIGINT,
	quantity BIGINT,
	role TEXT,
	reply TEXT,
	msg_id TEXT NOT NULL UNIQUE,
	PRIMARY KEY (guild_id, user_id),
	CONSTRAINT fk_guild_createitem FOREIGN KEY (guild_id)
        REFERENCES economy_config (guild_id) ON DELETE CASCADE
);
`,
}