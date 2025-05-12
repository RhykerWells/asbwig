package leaveserver

import (
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var Command = &dcommand.AsbwigCommand{
	Command:      "leaveserver",
	Category:	  dcommand.CategoryOwner,
	Description:  "Forces the bot to leave a given server",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		err := common.Session.GuildLeave(data.Args[0])
		if err == nil {
			common.Session.MessageReactionAdd(data.ChannelID, data.Message.ID, "👍")
		}
	}),
}
