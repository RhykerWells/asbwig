package set

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
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     []string{"set"},
	Description: "Changes the settings in the economy",
	Args: []*dcommand.Args {
		{Name: "Setting", Type: dcommand.String},
		{Name: "Value", Type: dcommand.String},
	},
	Run: settings,
}

func settings(data *dcommand.Data) {
	embed := &discordgo.MessageEmbed {Author: &discordgo.MessageEmbedAuthor{Name: data.Message.Author.Username, IconURL: data.Message.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: 0xFF0000}
	if len(data.Args) <= 0 {
		embed.Description = "No `settings` argument provided. Available arguments:\n`maxBet`, `startbalance`, `symbol` \nTo set it with the default settings use `default`"
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	setting := strings.ToLower(data.Args[0])
	o, _ := regexp.Compile("default|s(ymbol|tartbalance)|m(ax(bet)?|in)")
	if !o.MatchString(setting) {
		embed.Description = "Invalid `settings` argument provided. Available arguments:\n`maxBet`, `startbalance`, `symbol` \nTo set it with the default settings use `default`"
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	var settings models.EconomyConfig
	if setting == "default" {
		settings.GuildID = data.Message.GuildID
		settings.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("maxbet","symbol","startbalance",), boil.Infer())
		embed.Description = "Economy settings have been reset to default values"
		embed.Color = 0x00ff7b
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	if len(data.Args) <= 1 {
		embed.Description = "No `Value` argument provided"
		functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.Message.GuildID)).One(context.Background(), common.PQ)
	value := strings.ToLower(data.Args[1])
	switch setting {
	case "startbalance", "maxbet", "min", "max":
		nvalue := functions.ToInt64(value)
		if nvalue < 0 || (setting == "max" && nvalue == 0) {
			embed.Description = "Invalid `Value` argument provided"
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if setting == "max" && nvalue <= guild.Min {
			embed.Description = fmt.Sprintf("You can't set `max` to a value under `min`.\n`min` is currently set to %s%d", guild.Symbol, guild.Min)
			functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		displayvalue := fmt.Sprintf("%s%s", guild.Symbol, humanize.Comma(nvalue))
		if nvalue == 0 {
			displayvalue = "Disabled"
		}
		embed.Description = fmt.Sprintf("You set `%s` to %s", setting, displayvalue)
		embed.Color = 0x00ff7b
		value = fmt.Sprint(nvalue)
	case "symbol":
		embed.Description = fmt.Sprintf("You set `symbol` to %s", value)
		embed.Color = 0x00ff7b
	}
	query := fmt.Sprintf("UPDATE economy_config SET %s = $1 WHERE guild_id = $2", setting)
	queries.Raw(query, value, data.Message.GuildID).Exec(common.PQ)
	functions.SendMessage(data.Message.ChannelID, &discordgo.MessageSend{Embed: embed})
}