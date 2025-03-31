package viewsettings

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "viewsettings",
	Description: "Changes the settings in the economy",
	Run:         settings,
}

func settings(data *dcommand.Data) {
	guild, _ := common.Session.Guild(data.GuildID)
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " settings", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.SuccessGreen}
	guildConfig, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
	maxBet := ""
	symbol := guildConfig.Symbol
	startbalance := ""
	if guildConfig.Maxbet == 0 {
		maxBet = "Disabled"
	} else {
		maxBet = fmt.Sprint(symbol, humanize.Comma(guildConfig.Maxbet))
	}
	if guildConfig.Startbalance == 0 {
		startbalance = "Disabled"
	} else {
		startbalance = fmt.Sprint(symbol, humanize.Comma(guildConfig.Startbalance))
	}
	var workResponsesEnabled, crimeResponsesEnabled string = "Disabled", "Disabled"
	var workResponsesNum, crimeResponsesNum int = 0, 0
	if guildConfig.Customworkresponses {
		workResponsesEnabled = "Enabled"
	}
	if len(guildConfig.Workresponses) > 0 {
		workResponsesNum = len(guildConfig.Workresponses)
	}
	if guildConfig.Customcrimeresponses {
		crimeResponsesEnabled = "Enabled"
	}
	if len(guildConfig.Workresponses) > 0 {
		crimeResponsesNum = len(guildConfig.Crimeresponses)
	}
	min := fmt.Sprint(symbol, humanize.Comma(guildConfig.Min))
	max := fmt.Sprint(symbol, humanize.Comma(guildConfig.Max))
	embed.Description = fmt.Sprintf("min: `%s`\nmax: `%s`\nmaxBet: `%s`\nSymbol: `%s`\nstartBalance: `%s`\nWork responses: `%s` (%d)\nCrime responses: `%s` (%d)", min, max, maxBet, symbol, startbalance, workResponsesEnabled, workResponsesNum, crimeResponsesEnabled, crimeResponsesNum)
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
