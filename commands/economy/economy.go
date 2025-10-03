package economy

//go:generate sqlboiler --no-hooks psql

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/RhykerWells/asbwig/commands/economy/informational/balance"
	"github.com/RhykerWells/asbwig/commands/economy/informational/leaderboard"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/chickenFight"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/coinFlip"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/crime"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rob"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rollNumber"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/russianRoulette"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/snakeEyes"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/work"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/addMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/deposit"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/giveMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/removeMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/withdraw"
	"github.com/RhykerWells/asbwig/commands/economy/settings/addResponse"
	"github.com/RhykerWells/asbwig/commands/economy/settings/listResponses"
	"github.com/RhykerWells/asbwig/commands/economy/settings/removeResponse"
	"github.com/RhykerWells/asbwig/commands/economy/settings/set"
	"github.com/RhykerWells/asbwig/commands/economy/settings/viewsettings"
	"github.com/RhykerWells/asbwig/commands/economy/shop/buyItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/createItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/editItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/inventory"
	"github.com/RhykerWells/asbwig/commands/economy/shop/itemInfo"
	"github.com/RhykerWells/asbwig/commands/economy/shop/removeItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/sell"
	"github.com/RhykerWells/asbwig/commands/economy/shop/shop"
	"github.com/RhykerWells/asbwig/commands/economy/shop/useItem"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
)

// EconomySetup runs the following:
//  - The schema initialiser
//  - Registration of the guild and user join/leave functions
//  - Registration of the economy commands & their pagination
func EconomySetup(cmdHandler *dcommand.CommandHandler) {
	common.InitSchema("Economy", GuildEconomySchema...)
	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate){
		guildAddEconomyConfig,
	})
	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete){
		guildDeleteEconomyConfig,
	})
	events.RegisterGuildMemberJoinfunctions([]func(g *discordgo.GuildMemberAdd){
		guildMemberAddToEconomy,
	})
	events.RegisterGuildMemberLeavefunctions([]func(g *discordgo.GuildMemberRemove){
		guildMemberRemoveFromEconomy,
	})

	cmdHandler.RegisterCommands(
		//Info
		balance.Command,
		leaderboard.Command,
		//Moneymaking
		chickenfight.Command,
		coinflip.Command,
		crime.Command,
		rob.Command,
		rollnumber.Command,
		russianroulette.Command,
		snakeeyes.Command,
		work.Command,
		//Moneymanagement
		addmoney.Command,
		deposit.Command,
		givemoney.Command,
		removemoney.Command,
		withdraw.Command,
		//Settings
		addresponse.Command,
		listresponses.Command,
		removeresponse.Command,
		set.Command,
		viewsettings.Command,
		//Shop
		buyitem.Command,
		createitem.Command,
		edititem.Command,
		inventory.Command,
		iteminfo.Command,
		removeitem.Command,
		sell.Command,
		shop.Command,
		useItem.Command,
	)
	common.Session.AddHandler(leaderboard.Pagination)
	common.Session.AddHandler(listresponses.Pagination)
	common.Session.AddHandler(inventory.Pagination)
	common.Session.AddHandler(iteminfo.Pagination)
	common.Session.AddHandler(shop.Pagination)
}

// guildAddEconomyConfig creates the intial configs for the economy system for a specified guild
func guildAddEconomyConfig(g *discordgo.GuildCreate) {
	config := GetConfig(g.ID)
	SaveConfig(config)
}

// guildAddEconomyConfig deletes the configs for the economy system for a specified guild
func guildDeleteEconomyConfig(g *discordgo.GuildDelete) {
	config, err := models.EconomyConfigs(qm.Where("guild_id = ?", g.ID)).One(context.Background(), common.PQ)
	if err != nil {
		return
	}

	config.Delete(context.Background(), common.PQ)
}

// guildMemberAddToEconomy adds a member to the economy system
func guildMemberAddToEconomy(m *discordgo.GuildMemberAdd) {
	config := GetConfig(m.GuildID)
	userEntry := models.EconomyUser{
		GuildID: config.GuildID,
		UserID:  m.User.ID,
		Cash:    config.Startbalance,
		Bank:    0,
	}
	userEntry.Insert(context.Background(), common.PQ, boil.Infer())
}

// guildMemberRemoveFromEconomy removes a guild member from the economy system
func guildMemberRemoveFromEconomy(m *discordgo.GuildMemberRemove) {
	models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyUserInventories(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.User.ID)).DeleteAll(context.Background(), common.PQ)
}
