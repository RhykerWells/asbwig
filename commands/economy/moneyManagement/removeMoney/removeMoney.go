package removemoney

import (
	"context"
	"fmt"
	"strings"
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
	Command:     "removemoney",
	Description: "Removes money from a specified users cash/bank balance",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Place", Type: dcommand.String},
		{Name: "Amount", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		symbol := guild.Symbol
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
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
		if len(data.Args) <= 1 {
			embed.Description = "No `Destination` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		destination := strings.ToLower(data.Args[1])
		if destination != "cash" && destination != "bank" {
			embed.Description = "Invalid `Destination` argument provided\nPlease use `cash` or `bank`"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if len(data.Args) <= 2 {
			embed.Description = "No `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		amount := data.Args[2]
		if functions.ToInt64(amount) <= 0 {
			embed.Description = "Invalid `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if destination == "cash" {
			userCash, err := models.EconomyCashes(qm.Where("guild_id=? AND user_id=?", data.GuildID, member.User.ID)).One(context.Background(), common.PQ)
			var cash int64 = 0
			if err == nil {
				cash = userCash.Cash
			}
			cash = cash - functions.ToInt64(amount)
			cashEntry := models.EconomyCash{
				GuildID: data.GuildID,
				UserID: member.User.ID,
				Cash: cash,
			}
			_ = cashEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		} else {
			userBank, err := models.EconomyBanks(qm.Where("guild_id=? AND user_id=?", data.GuildID, member.User.ID)).One(context.Background(), common.PQ)
			var bank int64 = 0
			if err == nil {
				bank = userBank.Balance
			}
			bank = bank - functions.ToInt64(amount)
			bankEntry := models.EconomyBank{
				GuildID: data.GuildID,
				UserID: member.User.ID,
				Balance: bank,
			}
			_ = bankEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("balance"), boil.Infer())
		}
		embed.Description = fmt.Sprintf("You remove %s%s from %ss %s", symbol, humanize.Comma(functions.ToInt64(amount)), member.Mention(), destination)
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
