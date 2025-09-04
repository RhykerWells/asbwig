package viewsettings

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "viewsettings",
	Category:    dcommand.CategoryEconomy,
	Description: "Changes the settings in the economy",
	Run:         util.AdminOrManageServerCommand(func(data *dcommand.Data) { viewsettings(data) }),
}

func viewsettings(data *dcommand.Data) {
	guild, _ := common.Session.Guild(data.GuildID)
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " settings", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.SuccessGreen}
	guildConfig, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
	customWorkResponses, _ := models.EconomyCustomResponses(qm.Where("guild_id=? AND type='work'", data.GuildID)).All(context.Background(), common.PQ)
	customCrimeResponses, _ := models.EconomyCustomResponses(qm.Where("guild_id=? AND type='crime'", data.GuildID)).All(context.Background(), common.PQ)
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
	if len(customWorkResponses) > 0 {
		workResponsesEnabled = "Enabled"
		workResponsesNum = len(customWorkResponses)
	}
	if len(customCrimeResponses) > 0 {
		crimeResponsesEnabled = "Enabled"
		crimeResponsesNum = len(customCrimeResponses)
	}
	min := fmt.Sprint(symbol, humanize.Comma(guildConfig.Min))
	max := fmt.Sprint(symbol, humanize.Comma(guildConfig.Max))
	embed.Description = fmt.Sprintf("min: `%s`\nmax: `%s`\nmaxBet: `%s`\nSymbol: `%s`\nstartBalance: `%s`\nWork responses: `%s` (%d)\nCrime responses: `%s` (%d)", min, max, maxBet, symbol, startbalance, workResponsesEnabled, workResponsesNum, crimeResponsesEnabled, crimeResponsesNum)
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
