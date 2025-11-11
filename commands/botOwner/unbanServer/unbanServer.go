package unbanserver

import (
	"context"

	"github.com/RhykerWells/Summit/bot/core/models"
	"github.com/RhykerWells/Summit/bot/functions"
	"github.com/RhykerWells/Summit/commands/util"
	"github.com/RhykerWells/Summit/common"
	"github.com/RhykerWells/Summit/common/dcommand"
)

var Command = &dcommand.SummitCommand{
	Command:      "unbanserver",
	Category:     dcommand.CategoryOwner,
	Description:  "Removes the server ban from inviting the bot",
	ArgsRequired: 1,
	Args: []*dcommand.Arg{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		banned := util.IsGuildBanned(data.ParsedArgs[0].String())
		if !banned {
			functions.SendBasicMessage(data.ChannelID, "That guild was not banned")
		} else {
			models.BannedGuilds(models.BannedGuildWhere.GuildID.EQ(data.ParsedArgs[0].String())).DeleteAll(context.Background(), common.PQ)
		}
	}),
}
