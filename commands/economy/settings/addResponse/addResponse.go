package addresponse

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "addresponse",
	Category:    dcommand.CategoryEconomy,
	Description: "Adds a new response to use in `work` or `crime`",
	Args: []*dcommand.Args{
		{Name: "Type", Type: dcommand.ResponseType},
		{Name: "Response", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run:          util.AdminOrManageServerCommand(func(data *dcommand.Data) { addResponse(data) }),
}

func addResponse(data *dcommand.Data) {
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
	responseType := data.Args[0]
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
	responseEntry := models.EconomyCustomResponse{
		GuildID:  data.GuildID,
		Type:     responseType,
		Response: response,
	}
	responseEntry.Insert(context.Background(), common.PQ, boil.Infer())
	embed.Description = fmt.Sprintf("Successfully added `%s` to your list of responses", response)
	embed.Color = common.SuccessGreen
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
