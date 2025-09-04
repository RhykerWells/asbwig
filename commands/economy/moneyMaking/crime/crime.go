package crime

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "crime",
	Category:    dcommand.CategoryEconomy,
	Description: "Pew pew pew",
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=? AND type='crime'", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if cooldown.ExpiresAt.Time.After(time.Now()) {
				embed.Description = fmt.Sprintf("This command is on cooldown for <t:%d:R>", (time.Now().Unix() + int64(time.Until(cooldown.ExpiresAt.Time).Seconds())))
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		crimeResponses, _ := models.EconomyCustomResponses(qm.Where("guild_id=? AND type='crime'", data.GuildID)).All(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		amount := rand.Int63n(guild.Max - guild.Min)
		if rand.Int63n(2) == 1 {
			embed.Description = fmt.Sprintf("You broke the law for a pretty penny! You made %s%s in your crime spree", guild.Symbol, humanize.Comma(amount))
			if guild.Customcrimeresponses && len(crimeResponses) > 0 {
				embed.Description = strings.ReplaceAll(crimeResponses[rand.Intn(len(crimeResponses))].Response, "(amount)", fmt.Sprintf("%s%s", guild.Symbol, humanize.Comma(amount)))
			}
			embed.Color = common.SuccessGreen
			cash = cash + amount
		} else {
			embed.Description = fmt.Sprintf("You broke the law and got caught! You were arrested and lost %s%s", guild.Symbol, humanize.Comma(amount))
			cash = cash - amount
		}
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "crime", ExpiresAt: null.Time{Time: time.Now().Add(3600 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
