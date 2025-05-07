package coinflip

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
	Command:     "coinflip",
	Category: 	 dcommand.CategoryEconomy,
	Aliases:     []string{"cf", "flip"},
	Description: "Flips a coin. Head or tails. Payout is equal to `<Bet>`",
	Args: []*dcommand.Args{
		{Name: "Bet", Type: dcommand.Bet},
		{Name: "betSide", Type: dcommand.CoinSide},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		cooldown, err := models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=? AND type='coinflip'", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		if err == nil {
			if cooldown.ExpiresAt.Time.After(time.Now()) {
				embed.Description = fmt.Sprintf("This command is on cooldown for <t:%d:R>", (time.Now().Unix() + int64(time.Until(cooldown.ExpiresAt.Time).Seconds())))
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		amount := data.Args[0]
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
		betSide := data.Args[1]
		if rand.Int63n(2) == 1 {
			cash = cash + bet
			embed.Description = fmt.Sprintf("You flipped %s and won %s%s", betSide, guild.Symbol, humanize.Comma(bet))
			embed.Color = common.SuccessGreen
		} else {
			cash = cash + bet
			embed.Description = fmt.Sprintf("You flipped %s and lost", betSide)
		}
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "coinflip", ExpiresAt: null.Time{Time: time.Now().Add(300 * time.Second), Valid: true}}
		cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}