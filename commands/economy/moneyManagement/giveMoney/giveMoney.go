package givemoney

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "givemoney",
	Description: "Gives money to a specified users cash balance from your cash",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Amount", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		if len(data.Args) <= 0 {
			embed.Description = "No `User` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		receiving, err := functions.GetMember(data.GuildID, data.Args[0])
		if err != nil {
			embed.Description = "Invalid `User` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if len(data.Args) <= 1 {
			embed.Description = "No `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		amount := data.Args[1]
		if functions.ToInt64(amount) <= 0 {
			embed.Description = "Invalid `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		conversionAmount := functions.ToInt64(amount)
		if conversionAmount > cash {
			embed.Description = fmt.Sprintf("You don't have enough cash to give. You have %s%s", guild.Symbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		embed.Description = fmt.Sprintf("You gave %s%s to %s!", guild.Symbol, humanize.Comma(functions.ToInt64(amount)), receiving.Mention())
		embed.Color = common.SuccessGreen
		cash = cash - conversionAmount
		receivingUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var receivingCash int64 = 0
		if err == nil {
			receivingCash = receivingUser.Cash
		}
		receivingCash = receivingCash + conversionAmount
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		receivingEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: receivingCash}
		receivingEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
