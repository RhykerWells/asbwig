package economy

var GuildEconomySchema = []string{`
CREATE TABLE IF NOT EXISTS economy_config (
	guild_id TEXT PRIMARY KEY,
	min BIGINT NOT NULL DEFAULT '200',
	max BIGINT NOT NULL DEFAULT '500',
	maxbet BIGINT NOT NULL DEFAULT '5000',
	symbol TEXT NOT NULL DEFAULT '£',
	startbalance BIGINT NOT NULL DEFAULT '200',
	customworkresponses BOOL NOT NULL DEFAULT 'false',
	customcrimeresponses BOOL NOT NULL DEFAULT 'false',
	workresponses TEXT[],
	crimeresponses TEXT[]
);
`,`
CREATE TABLE IF NOT EXISTS economy_shop (
    guild_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price BIGINT NOT NULL CHECK (price > 0),
    quantity BIGINT,
    role TEXT,
    reply TEXT,
    soldby TEXT,
    UNIQUE (guild_id, name)
);
`,`
CREATE OR REPLACE FUNCTION deleteItemWhenOut()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity = 0 THEN
        DELETE FROM economy_shop WHERE guild_id = NEW.guild_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER trigger_delete_item_when_quantity_zero
AFTER UPDATE ON economy_shop
FOR EACH ROW
WHEN (NEW.quantity = 0)
EXECUTE FUNCTION deleteItemWhenOut();
`,`
CREATE TABLE IF NOT EXISTS economy_createitem (
	guild_id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	name TEXT,
	description TEXT,
	price BIGINT,
	quantity BIGINT,
	role TEXT,
	reply TEXT,
	expires_at TIMESTAMP,
	msg_id TEXT NOT NULL,
	UNIQUE (guild_id, user_id)
)
`,`
CREATE TABLE IF NOT EXISTS economy_cash (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	cash BIGINT NOT NULL,
	UNIQUE (guild_id, user_id)
);
`,`
CREATE TABLE IF NOT EXISTS economy_bank (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	balance BIGINT NOT NULL,
	UNIQUE (guild_id, user_id)
);
`,`
CREATE TABLE IF NOT EXISTS economy_cooldowns (
	ID SERIAL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	type TEXT NOT NULL,
	expires_at TIMESTAMP,
	UNIQUE (guild_id, user_id, type)
);`,
}