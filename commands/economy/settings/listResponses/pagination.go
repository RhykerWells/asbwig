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
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
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
	if b.MessageComponentData().CustomID == "responses_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	responseType := "type"
	if strings.Contains(embed[0].Author.Name, "crime-responses") {
		responseType = "crime"
	}
	offset := (page - 1) * 10
	display := ""
	guildResponses, err := models.EconomyCustomResponses(models.EconomyConfigWhere.GuildID.EQ(b.GuildID), models.EconomyCustomResponseWhere.Type.EQ(responseType), qm.Offset(offset)).All(context.Background(), common.PQ)
	if err != nil {
		display = "There are no responses on this page\nAdd some with `addresponse <Type> <Responses>`"
	}
	responseNumber := (page - 1) * 10
	for i, responses := range guildResponses {
		if i == 10 {
			break
		}
		responseNumber++
		display += fmt.Sprintf("%d) `%s`\n", responseNumber, responses.Response)
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
	if len(guildResponses) > responseNumber {
		row := components[0].(discordgo.ActionsRow)
		btnNext := row.Components[1].(discordgo.Button)
		btnNext.Disabled = false
		row.Components[1] = btnNext
		components[0] = row
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     embed,
			Components: components,
		},
	}
	common.Session.InteractionRespond(b.Interaction, response)
}
