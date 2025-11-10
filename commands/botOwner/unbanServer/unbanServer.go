package unbanserver

import (
	"context"

	"github.com/RhykerWells/summit/bot/core/models"
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/commands/util"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
)

var Command = &dcommand.SummitCommand{
	Command:      "unbanserver",
	Category:     dcommand.CategoryOwner,
	Description:  "Removes the server ban from inviting the bot",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
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
