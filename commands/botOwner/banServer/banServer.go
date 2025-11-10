package banserver

import (
	"context"

	"github.com/RhykerWells/summit/bot/core/models"
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/commands/util"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
)

var Command = &dcommand.SummitCommand{
	Command:      "ban",
	Category:     dcommand.CategoryOwner,
	Description:  "Bans a server from inviting the bot",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		banned := util.IsGuildBanned(data.ParsedArgs[0].String())
		if banned {
			functions.SendBasicMessage(data.ChannelID, "This guild is already banned")
		} else {
			guild := models.BannedGuild{
				GuildID: data.ParsedArgs[0].String(),
			}
			guild.Insert(context.Background(), common.PQ, boil.Infer())
		}
	}),
}
