package withdraw

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
	Command:     "withdraw",
	Description: "Withdraws a given amount from your bank",
	Args: []*dcommand.Args{
		{Name: "Amount", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash, bank int64
		if err == nil {
			cash = economyUser.Cash
			bank = economyUser.Bank
		}
		if len(data.Args) <= 0 {
			embed.Description = "No `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		value := data.Args[0]
		if functions.ToInt64(value) <= 0 && value != "all" {
			embed.Description = "Invalid `Amount` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		var amount int64
		if value == "all" {
			amount = bank
		} else {
			amount = functions.ToInt64(value)
		}
		if amount > bank {
			embed.Description = fmt.Sprintf("You're unable to withdraw more than you have in your bank\nYou currently have %s%s", guild.Symbol, humanize.Comma(bank))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if bank < 0 {
			embed.Description = fmt.Sprintf("You're unable to withdraw from your overdraft\nYou are currently %s%s in arrears", guild.Symbol, humanize.Comma(bank))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		cash = cash + amount
		bank = bank - amount
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash, Bank: bank}
		_ = userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash", "bank"), boil.Infer())
		embed.Description = fmt.Sprintf("You Withdrew %s%s from your bank\nThere is now %s%s in your bank", guild.Symbol, humanize.Comma(amount), guild.Symbol, humanize.Comma(bank))
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
