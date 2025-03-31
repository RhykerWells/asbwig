package listresponses

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Pagination(s *discordgo.Session, b *discordgo.InteractionCreate) {
	guild, _ := common.Session.Guild(b.GuildID)
	embed := []*discordgo.MessageEmbed{{Author: &discordgo.MessageEmbedAuthor{Name: b.Message.Embeds[0].Author.Name, IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "responses_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "response_forward"}}}}
	if b.MessageComponentData().CustomID != "responses_back" && b.MessageComponentData().CustomID != "responses_forward" {
		return
	}
	re := regexp.MustCompile(`\d+`)
	page, _ := strconv.Atoi(re.FindString(b.Message.Embeds[0].Footer.Text))
	if b.MessageComponentData().CustomID == "economy_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	responseType := "workresponses"
	if strings.Contains(embed[0].Author.Name, "crime-responses") {
		responseType = "crimeresponses"
	}
	offset :=  (page - 1) * 10
	display := ""
	guildConfig, err := models.EconomyConfigs(qm.Select("workresponses", "crimeresponses"), qm.Where("guild_id=?", b.GuildID), qm.Offset(offset)).One(context.Background(), common.PQ)
	var responses []string
	if err != nil {
		display = "There are no responses on this page"
	}  else {
		switch responseType {
		case "workresponses":
			responses = guildConfig.Workresponses
		case "crimeresponses":
			responses = guildConfig.Crimeresponses
		}
	}
	if len(responses) == 0 {
        display = "There are no responses on this page"
    } else {
        embed[0].Color = common.SuccessGreen
    }
	responseNumber := (page - 1) * 10
	for i, response := range responses {
		if i == 10 {
			break
		}
		responseNumber ++
		display += fmt.Sprintf("%d) `%s`\n", responseNumber, response)
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
	if page != 1 {
		row := components[0].(discordgo.ActionsRow)
		btnPrev := row.Components[0].(discordgo.Button)
		btnPrev.Disabled = false
		row.Components[0] = btnPrev
		components[0] = row	
	}
	if len(responses) > responseNumber {
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