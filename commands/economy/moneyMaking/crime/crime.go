package crime

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
	Command:     "crime",
	Description: "Pew pew pew",
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		payout := rand.Int63n(guild.Max - guild.Min)
		userCrime, err := models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=? AND type = 'crime'", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if userCrime.ExpiresAt.Time.After(time.Now()) {
				embed.Description = "This command is on cooldown"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		embed.Description = fmt.Sprintf("You commit a crime and stole all their money. You got %s%s from being a horrible person", guild.Symbol, humanize.Comma(payout))
		embed.Color = common.SuccessGreen
		userCash, err := models.EconomyCashes(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = userCash.Cash
		}
		userCash.Cash = cash + functions.ToInt64(payout)
		_, _ = userCash.Update(context.Background(), common.PQ, boil.Whitelist("cash"))
		cooldowns := models.EconomyCooldown{
			GuildID: data.GuildID,
			UserID:  data.Author.ID,
			Type:    "crime",
			ExpiresAt: null.Time{
				Time:  time.Now().Add(3600 * time.Second),
				Valid: true,
			},
		}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
