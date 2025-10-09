package rollnumber

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
	Command:     "rollnumber",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"roll", "rollnum"},
	Description: "Rolls a number\n**100** = payout of `<bet>*5`\n**90-99** = payout of `<Bet>*3`\n**65-89** = payout of `<Bet>`\n**64 and under** = Loss of `<Bet>`",
	Args: []*dcommand.Args{
		{Name: "Bet", Type: dcommand.Bet},
	},
	ArgsRequired: 1,
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(models.EconomyCooldownWhere.GuildID.EQ(data.GuildID), models.EconomyCooldownWhere.UserID.EQ(data.Author.ID), models.EconomyCooldownWhere.Type.EQ("rollnumber")).One(context.Background(), common.PQ)
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
			if guild.Maxbet > 0 {
				bet = guild.Maxbet
			} else {
				bet = cash
			}
		} else {
			bet = functions.ToInt64(amount)
		}
		if bet > cash {
			embed.Description = fmt.Sprintf("You can't bet more than you have in your hand. You currently have %s%s", guild.Symbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if guild.Maxbet > 0 && bet > guild.Maxbet {
			embed.Description = fmt.Sprintf("You can't bet more than the servers limit. The limit is %s%s", guild.Symbol, humanize.Comma(guild.Maxbet))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		roll := rand.Int63n(100) + 1
		embed.Color = common.SuccessGreen
		condition := "won"
		if roll >= 65 && roll < 90 {
			cash = cash + bet
		} else if roll >= 90 && roll < 100 {
			bet = bet * 3
			cash = cash + bet
		} else if roll == 100 {
			bet = bet * 5
			cash = cash + bet
		} else {
			cash = cash - bet
			condition = "lost"
			embed.Color = common.ErrorRed
		}
		embed.Description = fmt.Sprintf("The ball landed on %d, and you %s %s%s", roll, condition, guild.Symbol, humanize.Comma(bet))
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "rollnumber", ExpiresAt: null.Time{Time: time.Now().Add(300 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{models.EconomyCooldownColumns.GuildID, models.EconomyCooldownColumns.UserID, models.EconomyCooldownColumns.Type}, boil.Whitelist(models.EconomyCooldownColumns.ExpiresAt), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
