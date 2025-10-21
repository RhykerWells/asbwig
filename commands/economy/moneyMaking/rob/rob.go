package rob

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
	Command:     "rob",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"steal"},
	Description: "Pew pew pew",
	Args: []*dcommand.Args{
		{Name: "Member", Type: dcommand.Member},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(models.EconomyCooldownWhere.GuildID.EQ(data.GuildID), models.EconomyCooldownWhere.UserID.EQ(data.Author.ID), models.EconomyCooldownWhere.Type.EQ("rob")).One(context.Background(), common.PQ)
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
		if len(data.Args) <= 0 {
			embed.Description = "No `User` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		member, err := functions.GetMember(data.GuildID, data.Args[0])
		if err != nil {
			embed.Description = "Invalid `User` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if member.User.ID == data.Author.ID {
			embed.Description = "Invalid `User` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}

		victim, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(member.User.ID)).One(context.Background(), common.PQ)
		var victimCash int64 = 0
		if err == nil {
			victimCash = victim.Cash
		}
		if victimCash < 0 {
			embed.Description = "This user has no cash to steal :("
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		payout := rand.Int63n(victimCash) + 1
		embed.Description = fmt.Sprintf("You stole %s%s from %s", guild.EconomySymbol, humanize.Comma(payout), member.Mention())
		cash = cash + payout
		victimCash = victimCash - payout
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		victimEntry := models.EconomyUser{GuildID: data.GuildID, UserID: member.User.ID, Cash: victimCash}
		victimEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "work", ExpiresAt: null.Time{Time: time.Now().Add(18000 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{models.EconomyCooldownColumns.GuildID, models.EconomyCooldownColumns.UserID, models.EconomyCooldownColumns.Type}, boil.Whitelist(models.EconomyCooldownColumns.ExpiresAt), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
