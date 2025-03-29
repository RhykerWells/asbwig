package rob

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
	Command:     "rob",
	Description: "Pew pew pew",
	Args: []*dcommand.Args {
		{Name: "Member", Type: dcommand.User},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed {Author: &discordgo.MessageEmbedAuthor{Name: data.Message.Author.Username, IconURL: data.Message.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: 0xFF0000}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.Message.GuildID)).One(context.Background(), common.PQ)
		symbol := guild.Symbol
		userRob, err := models.EconomyCooldowns(qm.Where("guild_id = ? AND user_id = ? AND type = 'rob'", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if userRob.ExpiresAt.Time.After(time.Now()) {
				embed.Description = "This command is on cooldown"
				functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		if len(data.Args) <= 0 {
			embed.Description = "No `User` argument provided"
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		member, err := functions.GetMember(data.Message.GuildID, data.Args[0])
		if err != nil {
			embed.Description = "Invalid `User` argument provided"
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if member.User.ID == data.Message.Author.ID {
			embed.Description = "Invalid `User` argument provided"
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		victim, err := models.EconomyCashes(qm.Where("guild_id=? AND user_id=?", data.Message.GuildID, member.User.ID)).One(context.Background(), common.PQ)
		if err != nil || victim.Cash < 0 {
			embed.Description = "This user has no cash to steal :("
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		payout := rand.Int63n(victim.Cash)
		embed.Description = fmt.Sprintf("You stole %s%s from %s", symbol, humanize.Comma(payout), member.Mention())
		userCash, err := models.EconomyCashes(qm.Where("guild_id = ? AND user_id = ?", data.Message.GuildID, data.Message.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = userCash.Cash
		}
		userCash.Cash = cash + payout
		victim.Cash = victim.Cash - payout
		_, _ = victim.Update(context.Background(), common.PQ, boil.Whitelist("cash"))
		_, _ = userCash.Update(context.Background(), common.PQ, boil.Whitelist("cash"))
		cooldowns := models.EconomyCooldown{
			GuildID:  data.Message.GuildID,
			UserID:   data.Message.Author.ID,
			Type:     "work",
			ExpiresAt: null.Time{
				Time:  time.Now().Add(18000 * time.Second),
				Valid: true,
			},
		}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}