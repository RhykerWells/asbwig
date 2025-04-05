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
	Command: "leaderboard",
	Aliases: []string{"lb", "top"},
	Args: []*dcommand.Args{
		{Name: "Page", Type: dcommand.Int, Optional: true},
	},
	Description: "Views your server leaderboard",
	Run: (func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " leaderboard", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "economy_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "economy_forward"}}}}
		guildSettings, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		page := 1
		if len(data.Args) > 0 {
			page, _ = strconv.Atoi(data.Args[0])
			if page < 1 {
				page = 1
			}
		}
		offset := (page - 1) * 10
		display := ""
		economyUsers, err := models.EconomyUsers(qm.Where("guild_id=?", data.GuildID), qm.OrderBy("cash DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
		if err != nil || len(economyUsers) == 0 {
			display = "No users are in the leaderboard"
		} else {
			embed.Color = common.SuccessGreen
		}
		rank := (page - 1) * 10
		for i, entry := range economyUsers {
			if i == 10 {
				break
			}
			cash := humanize.Comma(entry.Cash)
			rank++
			drank := ""
			user, _ := functions.GetUser(entry.UserID)
			pos := map[int]string{1: "🥇", 2: "🥈", 3: "🥉"}
			_, exists := pos[rank]
			if exists {
				drank = pos[rank]
			} else {
				drank = fmt.Sprintf("  %d.", rank) // Default to number if no medal
			}
			display += fmt.Sprintf("**%v** %s **•** %s%s\n", drank, user.Username, guildSettings.Symbol, cash)
		}
		embed.Description = display
		embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
		if page != 1 {
			row := components[0].(discordgo.ActionsRow)
			btnPrev := row.Components[0].(discordgo.Button)
			btnPrev.Disabled = false
			row.Components[0] = btnPrev
			components[0] = row
		}
		if len(economyUsers) > rank {
			row := components[0].(discordgo.ActionsRow)
			btnNext := row.Components[1].(discordgo.Button)
			btnNext.Disabled = false
			row.Components[1] = btnNext
			components[0] = row
		}
		msg, _ := common.Session.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{Embed: embed, Components: components})
		go disableButtons(msg.ChannelID, msg.ID)
	}),
}

func disableButtons(channelID, messageID string) {
	time.Sleep(10 * time.Second)
	lbMessage, _ := common.Session.ChannelMessage(channelID, messageID)
	components := lbMessage.Components
	row := components[0].(*discordgo.ActionsRow)
	btnPrev := row.Components[0].(*discordgo.Button)
	btnNext := row.Components[1].(*discordgo.Button)
	btnPrev.Disabled = true
	btnNext.Disabled = true
	row.Components[0] = btnPrev
	row.Components[1] = btnNext
	components[0] = row
	message := &discordgo.MessageSend{
		Embed:      lbMessage.Embeds[0],
		Components: components,
	}
	functions.EditMessage(channelID, messageID, message)
}