package economy

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/RhykerWells/asbwig/commands/economy/informational/balance"
	"github.com/RhykerWells/asbwig/commands/economy/settings/set"
	"github.com/RhykerWells/asbwig/commands/economy/settings/viewsettings"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

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
)`}

func EconomySetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Economy", GuildEconomySchema...)

	cmdHandler.RegisterCommands(
		balance.Command,
		set.Command,
		viewsettings.Command,
	)
}

func GuildEconomyAdd(guild_id string) {
	const query = `
	SELECT guild_id FROM economy_config
	WHERE guild_id=$1
	`
	err := common.PQ.QueryRow(query, guild_id)
	if err != nil {
		guildEconomyDefault(guild_id)
	}
}

func guildEconomyDefault(guild_id string) {
	const query = `
	INSERT INTO economy_config (guild_id)
	VALUES ($1)
	`
	common.PQ.Exec(query, guild_id)
}