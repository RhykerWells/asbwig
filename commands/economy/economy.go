package economy

//go:generate sqlboiler --no-hooks psql

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/RhykerWells/asbwig/commands/economy/informational/balance"
	"github.com/RhykerWells/asbwig/commands/economy/informational/leaderboard"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	chickenfight "github.com/RhykerWells/asbwig/commands/economy/moneyMaking/chickenFight"
	coinflip "github.com/RhykerWells/asbwig/commands/economy/moneyMaking/coinFlip"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/crime"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rob"
	rollnumber "github.com/RhykerWells/asbwig/commands/economy/moneyMaking/rollNumber"
	russianroulette "github.com/RhykerWells/asbwig/commands/economy/moneyMaking/russianRoulette"
	snakeeyes "github.com/RhykerWells/asbwig/commands/economy/moneyMaking/snakeEyes"
	"github.com/RhykerWells/asbwig/commands/economy/moneyMaking/work"
	addmoney "github.com/RhykerWells/asbwig/commands/economy/moneyManagement/addMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/deposit"
	givemoney "github.com/RhykerWells/asbwig/commands/economy/moneyManagement/giveMoney"
	removemoney "github.com/RhykerWells/asbwig/commands/economy/moneyManagement/removeMoney"
	"github.com/RhykerWells/asbwig/commands/economy/moneyManagement/withdraw"
	addresponse "github.com/RhykerWells/asbwig/commands/economy/settings/addResponse"
	listresponses "github.com/RhykerWells/asbwig/commands/economy/settings/listResponses"
	removeresponse "github.com/RhykerWells/asbwig/commands/economy/settings/removeResponse"
	"github.com/RhykerWells/asbwig/commands/economy/settings/set"
	"github.com/RhykerWells/asbwig/commands/economy/settings/viewsettings"
	buyitem "github.com/RhykerWells/asbwig/commands/economy/shop/buyItem"
	createitem "github.com/RhykerWells/asbwig/commands/economy/shop/createItem"
	edititem "github.com/RhykerWells/asbwig/commands/economy/shop/editItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/inventory"
	iteminfo "github.com/RhykerWells/asbwig/commands/economy/shop/itemInfo"
	removeitem "github.com/RhykerWells/asbwig/commands/economy/shop/removeItem"
	"github.com/RhykerWells/asbwig/commands/economy/shop/sell"
	"github.com/RhykerWells/asbwig/commands/economy/shop/shop"
	"github.com/RhykerWells/asbwig/commands/economy/shop/useItem"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
)

// EconomySetup runs the following:
//   - The schema initialiser
//   - Registration of the guild and user join/leave functions
//   - Registration of the economy commands & their pagination
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
	config, err := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(g.ID)).One(context.Background(), common.PQ)
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
	models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(m.GuildID), models.EconomyUserWhere.UserID.EQ(m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyCooldowns(models.EconomyCooldownWhere.GuildID.EQ(m.GuildID), models.EconomyCooldownWhere.UserID.EQ(m.User.ID)).DeleteAll(context.Background(), common.PQ)
	models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(m.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(m.User.ID)).DeleteAll(context.Background(), common.PQ)
}
