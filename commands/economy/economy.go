package economy

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/RhykerWells/asbwig/commands/economy/informational/balance"
	set "github.com/RhykerWells/asbwig/commands/economy/informational/settings"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var GuildEconomySchema = []string{`
CREATE TABLE IF NOT EXISTS economy_config (
	guild_id TEXT PRIMARY KEY,
	max_bet BIGINT NOT NULL DEFAULT '5000',
	symbol TEXT NOT NULL DEFAULT '£',
	start_balance BIGINT NOT NULL DEFAULT '200'
);
`,`
CREATE TABLE IF NOT EXISTS economy_cash (
	guild_id TEXT PRIMARY KEY,
	user_id BIGINT NOT NULL,
	cash BIGINT
)
`}

func EconomySetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Economy", GuildEconomySchema...)

	cmdHandler.RegisterCommands(
		balance.Command,
		set.Command,
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