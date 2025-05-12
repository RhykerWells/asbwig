package unbanserver

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/RhykerWells/asbwig/common/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:      "unbanserver",
	Category:	  dcommand.CategoryOwner,
	Description:  "Removes the server ban from inviting the bot",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		banned := util.IsGuildBanned(data.Args[0])
		if !banned {
			functions.SendBasicMessage(data.ChannelID, "That guild was not banned")
		} else {
			models.BannedGuilds(qm.Where("guild_id=?", data.Args[0])).DeleteAll(context.Background(), common.PQ)
		}
	}),
}
