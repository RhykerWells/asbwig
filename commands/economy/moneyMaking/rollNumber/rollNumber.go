package rollnumber

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
	Command:     "rollnumber",
	Aliases: 	 []string{"roll", "num"},
	Description: "Rolls a number, 1-100 with a max chance of 5*<bet>",
	Args: []*dcommand.Args{
		{Name: "Bet", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		userRoll, err := models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=? AND type = 'rollnumber' WHERE expires_at > NOW()", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if userRoll.ExpiresAt.Time.After(time.Now()) {
				embed.Description = "This command is on cooldown"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		userCash, err := models.EconomyCashes(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = userCash.Cash
		}
		if len(data.Args) <= 0 {
			embed.Description = "No `Bet` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		amount := data.Args[0]
		if functions.ToInt64(amount) <= 0 && amount != "all" && amount != "max" {
			embed.Description = "Invalid `Bet` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		bet := int64(0)
		if amount == "all" {
			bet = cash
		} else if amount == "max" {
			if guild.Maxbet > 0 {
				bet = guild.Maxbet
			} else {
				bet = cash
			}
		} else {
			bet = functions.ToInt64(amount)
		}
		if bet > cash {
			embed.Description = fmt.Sprintf("You can't bet more than you have in your hand. You currently have %s%s", guild.Symbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if guild.Maxbet > 0 && bet > guild.Maxbet {
			embed.Description = fmt.Sprintf("You can't bet more than the servers limit. The limit is %s%s", guild.Symbol, humanize.Comma(guild.Maxbet))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		roll := rand.Int63n(100)
		embed.Color = common.SuccessGreen
		condition := "won"
		if roll >= 65 && roll < 90 {
			cash = cash + bet
		} else if roll >= 90 && roll < 100 {
			bet = bet * 3
			cash = cash + bet 
		} else if roll == 100 {
			bet = bet * 5
			cash = cash + bet
		} else {
			cash = cash - bet
			condition = "lost"
			embed.Color = common.ErrorRed
		}
		embed.Description = fmt.Sprintf("The ball landed on %d, and you %s %s%s", roll, condition, guild.Symbol, humanize.Comma(bet))
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		cashEntry := models.EconomyCash{
			GuildID: data.GuildID,
			UserID: data.Author.ID,
			Cash: cash,
		}
		_ = cashEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		cooldowns := models.EconomyCooldown{
			GuildID: data.GuildID,
			UserID:  data.Author.ID,
			Type:    "rollnumber",
			ExpiresAt: null.Time{
				Time:  time.Now().Add(300 * time.Second),
				Valid: true,
			},
		}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
	},
}
