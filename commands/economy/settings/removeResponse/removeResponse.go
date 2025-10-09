package removeresponse

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "removeresponse",
	Category:    dcommand.CategoryEconomy,
	Description: "Removes a response from being used in `work` or `crime`",
	Args: []*dcommand.Args{
		{Name: "Type", Type: dcommand.ResponseType},
		{Name: "Response", Type: dcommand.Int},
	},
	ArgsRequired: 2,
	Run:          util.AdminOrManageServerCommand(func(data *dcommand.Data) { removeResponse(data) }),
}

func removeResponse(data *dcommand.Data) {
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
	responseType := data.Args[0]
	responseToDelete := data.Args[1]
	if functions.ToInt64(responseToDelete) <= 0 {
		embed.Description = "Invalid `Response` argument provided. To view your responses. Use the `listresponses` command"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	guildResponses, err := models.EconomyCustomResponses(models.EconomyCustomResponseWhere.GuildID.EQ(data.GuildID), models.EconomyCustomResponseWhere.Type.EQ(responseType)).All(context.Background(), common.PQ)
	if err != nil {
		embed.Description = "There are no responses to delete."
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	var responseNumber int64 = 1
	for _, responses := range guildResponses {
		if responseNumber == functions.ToInt64(responseToDelete) {
			responseToDelete = responses.Response
			responses.Delete(context.Background(), common.PQ)
			continue
		}
		responseNumber++
	}
	embed.Description = fmt.Sprintf("Successfully removed `%s` from your list of responses", responseToDelete)
	embed.Color = common.SuccessGreen
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
