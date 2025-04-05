package removemoney

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
	Command:     "removemoney",
	Description: "Removes money from a specified users cash/bank balance",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Place", Type: dcommand.String},
		{Name: "Amount", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
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
		destination := data.Args[1]
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
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, member.User.ID)).One(context.Background(), common.PQ)
		var cash, bank int64
		if err == nil {
			cash = economyUser.Cash
			bank = economyUser.Bank
		}
		if destination == "cash" {
			cash = cash - functions.ToInt64(amount)
		} else {
			bank = bank - functions.ToInt64(amount)
		}
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: member.User.ID, Cash: cash, Bank: bank}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash", "bank"), boil.Infer())
		embed.Description = fmt.Sprintf("You removed %s%s from %ss %s", guild.Symbol, humanize.Comma(functions.ToInt64(amount)), member.Mention(), destination)
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
