package leaderboard

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "leaderboard",
	Aliases:     []string{"lb", "top"},
	Description: "Views your server leaderboard",
	Run: (func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " leaderboard", IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guildSettings, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		page := 1
		if len(data.Args) > 0 {
			page, _ = strconv.Atoi(data.Args[0])
			if page < 1 {
				page = 1
			}
		}
		offset :=  (page - 1) * 10
		guildCash, err := models.EconomyCashes(qm.Where("guild_id=?", data.GuildID), qm.OrderBy("cash DESC"), qm.Limit(10), qm.Offset(offset)).All(context.Background(), common.PQ)
		display := ""
		if err != nil || len(guildCash) == 0 {
			display = "No users are in the leaderboard"
		} else {
			embed.Color = common.SuccessGreen
		}
		rank := 0
		for i, entry := range guildCash {
			user, _ := functions.GetUser(entry.UserID)
			pos := map[int]string{1: "🥇", 2: "🥈", 3: "🥉"}
			rank = i + 1
			drank, exists := pos[rank]
			if !exists {
				drank = fmt.Sprintf("%d.", rank) // Default to number if no medal
			}
			display += fmt.Sprintf("**%v** %s **•** %s%s\n", drank, user.Username, guildSettings.Symbol, humanize.Comma(entry.Cash))
		}
		embed.Description = display
		embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}