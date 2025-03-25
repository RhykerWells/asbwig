package balance

import (
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
	Command:     []string{"balance"},
	Description: "Views your balance in the economy",
	Run: (func(data *dcommand.Data) {
		const query = `SELECT cash FROM economy_cash WHERE guild_id=$1 AND user_id=$2`
		var cash string
		err := common.PQ.QueryRow(query, data.Message.GuildID, data.Message.Author.ID).Scan(&cash)
		if err != nil {
			cash = "0"
		}
		embed := &discordgo.MessageEmbed {
			Author: &discordgo.MessageEmbedAuthor{
				Name:    data.Message.Author.Username,
				IconURL: data.Message.Author.AvatarURL("256"),
			},
			Description: fmt.Sprintf("%s's balance", data.Message.Author.Mention()),
			Fields: []*discordgo.MessageEmbedField {
				{
					Name: "cash",
					Value: cash,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Color: 0x00ff7b,
		}
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}
