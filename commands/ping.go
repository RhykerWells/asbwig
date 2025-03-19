package commands

import (
	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/bwmarrin/discordgo"
)

var Ping = &AsbwigCommand {
	Command:	[]string{"ping"},
	Description: "Displays bot latency",
	Run: (func(data *Data) {
		message := &discordgo.MessageSend {
			Content: "Weee",
		}
		functions.SendMessage(data.Message.ChannelID, message)
	}),
}