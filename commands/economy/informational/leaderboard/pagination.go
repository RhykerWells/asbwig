package leaderboard

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Pagination(s *discordgo.Session, b *discordgo.InteractionCreate) {
	guild, _ := common.Session.Guild(b.GuildID)
	embed := []*discordgo.MessageEmbed{{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " leaderboard", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "economy_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "economy_forward"}}}}
	guildSettings, _ := models.EconomyConfigs(qm.Where("guild_id=?", b.GuildID)).One(context.Background(), common.PQ)
	if b.MessageComponentData().CustomID != "economy_back" && b.MessageComponentData().CustomID != "economy_forward" {
		return
	}
	re := regexp.MustCompile(`\d+`)
	page, _ := strconv.Atoi(re.FindString(b.Message.Embeds[0].Footer.Text))
	if b.MessageComponentData().CustomID == "economy_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	offset :=  (page - 1) * 10
	guildCash, err := models.EconomyCashes(qm.Where("guild_id=?", b.GuildID), qm.OrderBy("cash DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
	display := ""
	if err != nil || len(guildCash) == 0 {
		display = "No users are in the leaderboard"
	} else {
		embed[0].Color = common.SuccessGreen
	}
	rank := (page - 1) * 10
	for i, entry := range guildCash {
		if i == 10 {
			break
		}
		cash := humanize.Comma(entry.Cash)
		rank ++
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
	embed[0].Description = display
	embed[0].Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
	if page != 1 {
		row := components[0].(discordgo.ActionsRow)
		btnPrev := row.Components[0].(discordgo.Button)
		btnPrev.Disabled = false
		row.Components[0] = btnPrev
		components[0] = row	
	}
	if len(guildCash) > rank {
		row := components[0].(discordgo.ActionsRow)
		btnNext := row.Components[1].(discordgo.Button)
		btnNext.Disabled = false
		row.Components[1] = btnNext
		components[0] = row		
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
			Components: components,
		},
	}
	common.Session.InteractionRespond(b.Interaction, response)
}