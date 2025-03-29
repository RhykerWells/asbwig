package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "balance",
	Description: "Views your balance in the economy",
	Run: (func(data *dcommand.Data) {
		userCash, err := models.EconomyCashes(qm.Where("guild_id = ? AND user_id = ?", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = userCash.Cash
		}
		userBank, err := models.EconomyBanks(qm.Where("guild_id = ? AND user_id = ?", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		var bank int64 = 0
		if err == nil {
			bank = userBank.Balance
		}
		networth := cash + bank
		embed := &discordgo.MessageEmbed {
			Author: &discordgo.MessageEmbedAuthor{
				Name:    data.Message.Author.Username,
				IconURL: data.Message.Author.AvatarURL("256"),
			},
			Description: fmt.Sprintf("%s's balance", data.Message.Author.Mention()),
			Fields: []*discordgo.MessageEmbedField {
				{
					Name: "Cash",
					Value: fmt.Sprint(cash),
					Inline: true,
				},
				{
					Name: "Bank",
					Value: fmt.Sprint(bank),
					Inline: true,
				},
				{
					Name: "Networth",
					Value: fmt.Sprint(networth),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Color: 0x00ff7b,
		}
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}
