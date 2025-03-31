package economy

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/RhykerWells/asbwig/commands/economy/informational/balance"
	"github.com/RhykerWells/asbwig/commands/economy/informational/leaderboard"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/coinFlip"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/crime"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rob"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rollNumber"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/snakeEyes"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/work"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/addMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/deposit"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/removeMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/giveMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/withdraw"
	"github.com/RhykerWells/asbwig/commands/economy/settings/set"
	"github.com/RhykerWells/asbwig/commands/economy/settings/viewsettings"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

func EconomySetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Economy", GuildEconomySchema...)

	cmdHandler.RegisterCommands(
		balance.Command,
		leaderboard.Command,
		work.Command,
		crime.Command,
		rob.Command,
		rollnumber.Command,
		snakeeyes.Command,
		coinflip.Command,
		addmoney.Command,
		removemoney.Command,
		deposit.Command,
		withdraw.Command,
		givemoney.Command,
		set.Command,
		viewsettings.Command,
	)
	common.Session.AddHandler(leaderboard.Pagination)
}

func GuildEconomyAdd(guild_id string) {
	const query = `SELECT guild_id FROM economy_config WHERE guild_id=$1`
	err := common.PQ.QueryRow(query, guild_id)
	if err != nil {
		guildEconomyDefault(guild_id)
	}
}

func guildEconomyDefault(guild_id string) {
	const query = `INSERT INTO economy_config (guild_id) VALUES ($1)`
	common.PQ.Exec(query, guild_id)
}