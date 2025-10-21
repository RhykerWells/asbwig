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
	guildConfig, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)

	var workResponsesEnabled, crimeResponsesEnabled bool = guildConfig.EconomyCustomWorkResponsesEnabled, guildConfig.EconomyCustomCrimeResponsesEnabled
	var workResponsesNum, crimeResponsesNum int = len(guildConfig.EconomyCustomWorkResponses), len(guildConfig.EconomyCustomCrimeResponses)
	
	maxBet := ""
	symbol := guildConfig.EconomySymbol
	startbalance := ""
	if guildConfig.EconomyMaxBet == 0 {
		maxBet = "Disabled"
	} else {
		maxBet = fmt.Sprint(symbol, humanize.Comma(guildConfig.EconomyMaxBet))
	}
	if guildConfig.EconomyStartBalance == 0 {
		startbalance = "Disabled"
	} else {
		startbalance = fmt.Sprint(symbol, humanize.Comma(guildConfig.EconomyStartBalance))
	}

	min := fmt.Sprint(symbol, humanize.Comma(guildConfig.EconomyMinReturn))
	max := fmt.Sprint(symbol, humanize.Comma(guildConfig.EconomyMaxReturn))
	embed.Description = fmt.Sprintf("min: `%s`\nmax: `%s`\nmaxBet: `%s`\nSymbol: `%s`\nstartBalance: `%s`\nWork responses: `%s` (%d)\nCrime responses: `%s` (%d)", min, max, maxBet, symbol, startbalance, workResponsesEnabled, workResponsesNum, crimeResponsesEnabled, crimeResponsesNum)
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
