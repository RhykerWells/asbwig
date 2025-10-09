package removemoney

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "removemoney",
	Category:    dcommand.CategoryEconomy,
	Description: "Removes money from a specified users cash/bank balance",
	Args: []*dcommand.Args{
		{Name: "Member", Type: dcommand.Member},
		{Name: "Place", Type: dcommand.UserBalance},
		{Name: "Amount", Type: dcommand.Int},
	},
	ArgsRequired: 3,
	Run: util.AdminOrManageServerCommand(func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		member, _ := functions.GetMember(data.GuildID, data.Args[0])
		destination := data.Args[1]
		amount := data.Args[2]
		economyUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(member.User.ID)).One(context.Background(), common.PQ)
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
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash, models.EconomyUserColumns.Bank), boil.Infer())
		embed.Description = fmt.Sprintf("You removed %s%s from %ss %s", guild.Symbol, humanize.Comma(functions.ToInt64(amount)), member.Mention(), destination)
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}
