package work

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "work",
	Description: "Work work work",
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed {Author: &discordgo.MessageEmbedAuthor{Name: data.Message.Author.Username, IconURL: data.Message.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: 0xFF0000}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.Message.GuildID)).One(context.Background(), common.PQ)
		payout := rand.Int63n(guild.Max - guild.Min)
		userWork, err := models.EconomyCooldowns(qm.Where("guild_id = ? AND user_id = ? AND type = 'work'", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if userWork.ExpiresAt.Time.After(time.Now()) {
				embed.Description = "This command is on cooldown"
				functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		embed.Description = fmt.Sprintf("You decided to work today! You got paid a hefty %s%s", guild.Symbol, humanize.Comma(payout))
		embed.Color = 0x00ff7b
		userCash, err := models.EconomyCashes(qm.Where("guild_id = ? AND user_id = ?", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = userCash.Cash
		}
		userCash.Cash = cash + functions.ToInt64(payout)
		_, _ = userCash.Update(context.Background(), common.PQ, boil.Whitelist("cash"))
		cooldowns := models.EconomyCooldown{
			GuildID:  data.Message.GuildID,
			UserID:   data.Message.Author.ID,
			Type:     "work",
			ExpiresAt: null.Time{
				Time:  time.Now().Add(3600 * time.Second),
				Valid: true,
			},
		}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}