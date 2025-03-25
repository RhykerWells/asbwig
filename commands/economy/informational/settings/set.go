package set

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var Command = &dcommand.AsbwigCommand{
	Command:     []string{"set"},
	Description: "Changes the settings in the economy",
	Run: (func(data *dcommand.Data) {
		opt := "\nAvailable settings: `betMax`, `startbalance`, `symbol` \nTo set it with the default settings use `default`"
		var settings models.EconomyConfig
		embed := &discordgo.MessageEmbed {
			Author: &discordgo.MessageEmbedAuthor{
				Name:    data.Message.Author.Username,
				IconURL: data.Message.Author.AvatarURL("256"),
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Color: 0xFF0000,
		}
		if len(data.Args) <= 0 {
			embed.Description = fmt.Sprintf("No `settings` argument provided. Options are %s", opt)
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		setting := data.Args[0]
		if setting == "default" {
			settings.GuildID = data.Message.GuildID
			settings.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("max_bet", "symbol", "start_balance"), boil.Infer())
			embed.Description = "Economy settings have been reset to default values"
			embed.Color = 0x00ff7b
		}
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
	}),
}
