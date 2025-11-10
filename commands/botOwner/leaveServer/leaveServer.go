package leaveserver

import (
	"github.com/RhykerWells/summit/commands/util"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
)

var Command = &dcommand.SummitCommand{
	Command:      "leaveserver",
	Category:     dcommand.CategoryOwner,
	Description:  "Forces the bot to leave a given server",
	ArgsRequired: 1,
	Args: []*dcommand.Arg{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		err := common.Session.GuildLeave(data.ParsedArgs[0].String())
		if err == nil {
			common.Session.MessageReactionAdd(data.ChannelID, data.Message.ID, "👍")
		}
	}),
}
