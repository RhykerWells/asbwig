package setstatus

import (
	"strings"

	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/Ranger-4297/asbwig/commands/util"
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand {
	Command:		[]string{"setstatus"},
	Description: 	"Changes the bot status",
	ArgsRequired:	1,
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		functions.SetStatus(strings.Join(data.Args, " "))
		message := &discordgo.MessageSend {
			Content: "Status changed",
		}
		functions.SendMessage(data.Message.ChannelID, message)
	}),
}