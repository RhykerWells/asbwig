package balance

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "balance",
	Aliases:     []string{"bal"},
	Description: "Views your balance in the economy",
	Run: (func(data *dcommand.Data) {
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash, bank int64 = 0, 0
		if err == nil {
			cash = economyUser.Cash
			bank = economyUser.Bank
		}
		rankQuery := `SELECT position FROM (SELECT user_id, RANK() OVER (ORDER BY cash DESC) AS position FROM economy_cashs WHERE guild_id=$1) AS s WHERE user_id=$2`
		var rank int64
		drank := ""
		row := common.PQ.QueryRow(rankQuery, data.GuildID, data.Author.ID).Scan(&rank)
		if row == sql.ErrNoRows {
			rank = 0
		}
		if rank > 0 {
			ord := "th"
			cent := functions.ToInt64(math.Mod(float64(rank), float64(100)))
			dec := functions.ToInt64(math.Mod(float64(rank), float64(10)))
			if cent < int64(10) || cent > int64(19) {
				logrus.Println("cent is NOT between 10 and 19 2")
				if dec == 1 {
					ord = "st"
				} else if dec == 2 {
					ord = "nd"
				} else if dec == 3 {
					ord = "rd"
				}
			}
			drank = fmt.Sprintf("%d%s.", rank, ord)
		} else {
			drank = "None"
		}
		networth := cash + bank
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    data.Author.Username,
				IconURL: data.Author.AvatarURL("256"),
			},
			Description: fmt.Sprintf("%s's balance\nLeaderboard rank %s", data.Author.Mention(), drank),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Cash",
					Value:  fmt.Sprint(cash),
					Inline: true,
				},
				{
					Name:   "Bank",
					Value:  fmt.Sprint(bank),
					Inline: true,
				},
				{
					Name:   "Networth",
					Value:  fmt.Sprint(networth),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     common.SuccessGreen,
		}
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}
