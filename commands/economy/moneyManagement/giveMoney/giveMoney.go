package givemoney

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "givemoney",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"loan"},
	Description: "Gives money to a specified users cash balance from your cash",
	Args: []*dcommand.Args{
		{Name: "Member", Type: dcommand.Member},
		{Name: "Amount", Type: dcommand.Int},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		receiving, _ := functions.GetMember(data.GuildID, data.Args[0])
		amount := data.Args[1]
		conversionAmount := functions.ToInt64(amount)
		if conversionAmount > cash {
			embed.Description = fmt.Sprintf("You don't have enough cash to give. You have %s%s", guild.EconomySymbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		embed.Description = fmt.Sprintf("You gave %s%s to %s!", guild.EconomySymbol, humanize.Comma(functions.ToInt64(amount)), receiving.Mention())
		embed.Color = common.SuccessGreen
		cash = cash - conversionAmount
		receivingUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(receiving.User.ID)).One(context.Background(), common.PQ)
		var receivingCash int64 = 0
		if err == nil {
			receivingCash = receivingUser.Cash
		}
		receivingCash = receivingCash + conversionAmount
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		receivingEntry := models.EconomyUser{GuildID: data.GuildID, UserID: receiving.User.ID, Cash: receivingCash}
		receivingEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
