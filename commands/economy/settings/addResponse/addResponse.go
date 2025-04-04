package addresponse

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "addresponse",
	Description: "Adds a new response to use in `work` or `crime`",
	Args: []*dcommand.Args{
		{Name: "Type", Type: dcommand.String},
		{Name: "Response", Type: dcommand.String},
	},
	Run: addResponse,
}

func addResponse(data *dcommand.Data) {
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
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
	if len(data.Args) <= 1 {
		embed.Description = "No `Response` argument provided. Please include the exact string `(amount)` as a placeholder for where the amount goes"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return 
	}
	response := strings.Join(data.ArgsNotLowered[1:], " ")
	for _, char := range "\"" {
		response = strings.ReplaceAll(response, string(char), "")
	}
	a, _ := regexp.Compile(`\(amount\)`)
	if !a.MatchString(response) {
		embed.Description = "Invalid `Response` argument provided. Please include the exact string `(amount)` as a placeholder for where the amount goes"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return 
	}
	responseEntry := models.EconomyCustomResponse {
		GuildID: data.GuildID,
		Type: responseType,
		Response: response,
	}
	_ = responseEntry.Insert(context.Background(), common.PQ, boil.Infer())
	embed.Description = fmt.Sprintf("Successfully added `%s` to your list of responses", response)
	embed.Color = common.SuccessGreen
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}