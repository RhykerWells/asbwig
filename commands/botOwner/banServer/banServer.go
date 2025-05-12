package banserver

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/RhykerWells/asbwig/common/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var Command = &dcommand.AsbwigCommand{
	Command:      "ban",
	Category:	  dcommand.CategoryOwner,
	Description:  "Bans a server from inviting the bot",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		banned := util.IsGuildBanned(data.Args[0])
		if banned {
			functions.SendBasicMessage(data.ChannelID, "This guild is already banned")
		} else {
			guild := models.BannedGuild{
				GuildID: data.Args[0],
			}
			guild.Insert(context.Background(), common.PQ, boil.Infer())
		}
	}),
}