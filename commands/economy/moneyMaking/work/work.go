package work

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
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
	Command:     "work",
	Category:    dcommand.CategoryEconomy,
	Description: "Work work work",
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(models.EconomyCooldownWhere.GuildID.EQ(data.GuildID), models.EconomyCooldownWhere.UserID.EQ(data.Author.ID), models.EconomyCooldownWhere.Type.EQ("work")).One(context.Background(), common.PQ)
		if err == nil {
			if cooldown.ExpiresAt.Time.After(time.Now()) {
				embed.Description = fmt.Sprintf("This command is on cooldown for <t:%d:R>", (time.Now().Unix() + int64(time.Until(cooldown.ExpiresAt.Time).Seconds())))
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		guild, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		workResponses := guild.EconomyCustomWorkResponses
		economyUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		payout := rand.Int63n(guild.EconomyMaxReturn - guild.EconomyMinReturn)
		embed.Description = fmt.Sprintf("You decided to work today! You got paid a hefty %s%s", guild.EconomySymbol, humanize.Comma(payout))
		if guild.EconomyCustomWorkResponsesEnabled && len(workResponses) > 0 {
			embed.Description = strings.ReplaceAll(workResponses[rand.Intn(len(workResponses))], "(amount)", fmt.Sprintf("%s%s", guild.EconomySymbol, humanize.Comma(payout)))
		}
		embed.Color = common.SuccessGreen
		cash = cash + functions.ToInt64(payout)
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "work", ExpiresAt: null.Time{Time: time.Now().Add(3600 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{models.EconomyCooldownColumns.GuildID, models.EconomyCooldownColumns.UserID, models.EconomyCooldownColumns.Type}, boil.Whitelist(models.EconomyCooldownColumns.ExpiresAt), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
