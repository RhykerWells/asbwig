package listresponses

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "listresponses",
	Category: 	 dcommand.CategoryEconomy,
	Description: "Lists all responses for  `work` or `crime`",
	Args: []*dcommand.Args{
		{Name: "Type", Type: dcommand.String},
	},
	Run: util.AdminOrManageServerCommand(func(data *dcommand.Data) {listResponses(data)}),
}

func listResponses(data *dcommand.Data) {
	guild, _ := common.Session.Guild(data.GuildID)
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + "responses", IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "responses_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "responses_forward"}}}}
	if len(data.Args) <= 0 {
		embed.Description = "No `Type` argument provided. Available arguments:\n`Work`, `Crime`"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	responseType := data.Args[0]
	if responseType != "work" && responseType != "crime" {
		embed.Description = "Invalid `Type` argument provided. Available arguments:\n`Work`, `Crime`"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	embedAuthor := fmt.Sprintf("%s %s-responses", guild.Name, responseType)
	page := 1
	if len(data.Args) > 1 {
		page, _ = strconv.Atoi(data.Args[1])
		if page < 1 {
			page = 1
		}
	}
	offset := (page - 1) * 10
	display := ""
	guildResponses, _ := models.EconomyCustomResponses(qm.Where("guild_id=? AND type=?", data.GuildID, responseType), qm.Offset(offset)).All(context.Background(), common.PQ)
	if len(guildResponses) <= 0 {
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
	embed.Author = &discordgo.MessageEmbedAuthor{Name: embedAuthor, IconURL: guild.IconURL("256")}
	embed.Description = display
	embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
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
	msg, _ := common.Session.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{Embed: embed, Components: components})
	go disableButtons(msg.ChannelID, msg.ID)
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