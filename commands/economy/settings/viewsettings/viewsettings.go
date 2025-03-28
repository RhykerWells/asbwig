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
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.SuccessGreen}
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
	maxBet := ""
	symbol := guild.Symbol
	startbalance := ""
	if guild.Maxbet == 0 {
		maxBet = "Disabled"
	} else {
		maxBet = fmt.Sprint(symbol, humanize.Comma(guild.Maxbet))
	}
	if guild.Startbalance == 0 {
		startbalance = "Disabled"
	} else {
		startbalance = fmt.Sprint(symbol, humanize.Comma(guild.Startbalance))
	}
	min := fmt.Sprint(symbol, humanize.Comma(guild.Min))
	max := fmt.Sprint(symbol, humanize.Comma(guild.Max))
	embed.Description = fmt.Sprintf("min: `%s`\nmax: `%s`\nmaxBet: `%s`\nSymbol: `%s`\nstartBalance: `%s`", min, max, maxBet, symbol, startbalance)
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}
