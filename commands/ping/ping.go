package ping

import (
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand {
	Command:	[]string{"ping"},
	Description: "Displays bot latency",
	Run: (func(data *dcommand.Data) {
		message := &discordgo.MessageSend {
			Content: "Weee",
		}
		functions.SendMessage(data.Message.ChannelID, message)
	}),
}