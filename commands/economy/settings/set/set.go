package set

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "set",
	Category: 	 dcommand.CategoryEconomy,
	Description: "Changes the settings in the economy",
	Args: []*dcommand.Args{
		{Name: "Setting", Type: dcommand.String},
		{Name: "Value", Type: dcommand.String},
	},
	Run: util.AdminOrManageServerCommand(func(data *dcommand.Data) {settings(data)}),
}

func settings(data *dcommand.Data) {
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
	if len(data.Args) <= 0 {
		embed.Description = "No `settings` argument provided. Available arguments:\n`maxBet`, `startBalance`, `symbol` `workResponses`, `crimeResponses`\nTo set it with the default settings use `default`"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	setting := data.Args[0]
	o, _ := regexp.Compile("default|s(ymbol|tartbalance)|m(ax(bet)?|in)|(work|crime)responses")
	if !o.MatchString(setting) {
		embed.Description = "Invalid `settings` argument provided. Available arguments:\n`maxBet`, `startbalance`, `symbol` \nTo set it with the default settings use `default`"
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	var settings models.EconomyConfig
	if setting == "default" {
		settings.GuildID = data.GuildID
		settings.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("maxbet", "symbol", "startbalance"), boil.Infer())
		embed.Description = "Economy settings have been reset to default values"
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	var arg string
	switch setting {
	case "startbalance", "maxbet", "min", "max":
		arg = "<Value:Amount>"
	case "symbol":
		arg = "<Value:Emoji/CurrencySymbol>"
	case "workresponses", "crimeresponses":
		arg = "<Value:Enabled/Disable"
	}
	if len(data.Args) <= 1 {
		embed.Description = fmt.Sprintf("No `Value` argument provided for `%s`. Available arguments:\n`%s`", setting, arg)
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
		return
	}
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
	value := data.Args[1]
	switch setting {
	case "startbalance", "maxbet", "min", "max":
		nvalue := functions.ToInt64(value)
		if nvalue < 0 || (setting == "max" && nvalue == 0) {
			embed.Description = fmt.Sprintf("Invalid `Value` argument provided for %s. Available arguments:\n`%s`", setting, arg)
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if setting == "max" && nvalue <= guild.Min {
			embed.Description = fmt.Sprintf("You can't set `max` to a value under `min`.\n`min` is currently set to %s%d", guild.Symbol, guild.Min)
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		displayvalue := fmt.Sprintf("%s%s", guild.Symbol, humanize.Comma(nvalue))
		if nvalue == 0 {
			displayvalue = "Disabled"
		}
		embed.Description = fmt.Sprintf("You set `%s` to %s", setting, displayvalue)
		embed.Color = common.SuccessGreen
		value = fmt.Sprint(nvalue)
	case "symbol":
		embed.Description = fmt.Sprintf("You set `symbol` to %s", value)
		embed.Color = common.SuccessGreen
	case "workresponses", "crimeresponses":
		if value != "enabled" && value != "disabled" {
			embed.Description = fmt.Sprintf("Invalid `Value` argument provided for %s. Available arguments:\n`%s`", setting, arg)
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		embed.Description = fmt.Sprintf("You set `%s` to `%s`", setting, value)
		embed.Color = common.SuccessGreen
		switch setting {
		case "workresponses":
			setting = "customworkresponses"
		case "crimeresponses":
			setting = "customcrimeresponses"
		}
		switch value {
		case "enabled":
			value = "true"
		case "disabled":
			value = "false"
		}
	}
	query := fmt.Sprintf("UPDATE economy_config SET %s=$1 WHERE guild_id=$2", setting)
	queries.Raw(query, value, data.GuildID).Exec(common.PQ)
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}