package chickenfight

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "chickenfight",
	Category:    dcommand.CategoryEconomy,
	Description: "Chicken fight for a payout of <Bet> with a base payout of 50%. Increases each win up to 70%",
	Args: []*dcommand.Args{
		{Name: "Bet", Type: dcommand.Bet},
	},
	ArgsRequired: 1,
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(models.EconomyCooldownWhere.GuildID.EQ(data.GuildID), models.EconomyCooldownWhere.UserID.EQ(data.Author.ID), models.EconomyCooldownWhere.Type.EQ("chickenfight")).One(context.Background(), common.PQ)
		if err == nil {
			if cooldown.ExpiresAt.Time.After(time.Now()) {
				embed.Description = fmt.Sprintf("This command is on cooldown for <t:%d:R>", (time.Now().Unix() + int64(time.Until(cooldown.ExpiresAt.Time).Seconds())))
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		guild, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		amount := data.Args[0]
		bet := int64(0)
		if amount == "all" {
			bet = cash
		} else if amount == "max" {
			if guild.EconomyMaxBet > 0 {
				bet = guild.EconomyMaxBet
			} else {
				bet = cash
			}
		} else {
			bet = functions.ToInt64(amount)
		}
		if bet > cash {
			embed.Description = fmt.Sprintf("You can't bet more than you have in your hand. You currently have %s%s", guild.EconomySymbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if guild.EconomyMaxBet > 0 && bet > guild.EconomyMaxBet {
			embed.Description = fmt.Sprintf("You can't bet more than the servers limit. The limit is %s%s", guild.EconomySymbol, humanize.Comma(guild.EconomyMaxBet))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		chicken, exists := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(data.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(data.Author.ID), models.EconomyUserInventoryWhere.Name.IN([]string{"Chicken", "chicken"})).One(context.Background(), common.PQ)
		shopChicken := "chicken"
		shopChickenItem, ok := models.EconomyShops(models.EconomyShopWhere.GuildID.EQ(data.GuildID), models.EconomyShopWhere.Name.IN([]string{"Chicken", "chicken"})).One(context.Background(), common.PQ)
		if ok == nil {
			shopChicken = shopChickenItem.Name
		}
		if exists != nil {
			embed.Description = fmt.Sprintf("You don't have this item\nBuy it with `buyitem %s`", shopChicken)
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		win := false
		winChance := economyUser.Cfwinchance
		chance := rand.Int63n(100) + 1
		if chance <= winChance {
			win = true
		}
		if win {
			if winChance != 70 {
				winChance = winChance + 1
			}
			cash = cash + bet
			embed.Description = "Your chicken won the fight! Play again with an increased win chance"
			embed.Color = common.SuccessGreen
		} else {
			winChance = 50
			cash = cash - bet
			embed.Description = "Your chicken lost the fight and died :(\nBuy a new one to play again"
			chicken.Delete(context.Background(), common.PQ)
		}
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash, Cfwinchance: winChance}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash, models.EconomyUserColumns.Cfwinchance), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "chickenfight", ExpiresAt: null.Time{Time: time.Now().Add(300 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{models.EconomyCooldownColumns.GuildID, models.EconomyCooldownColumns.UserID, models.EconomyCooldownColumns.Type}, boil.Whitelist(models.EconomyCooldownColumns.ExpiresAt), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
