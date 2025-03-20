package help

import (
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand {
	Command:	[]string{"help"},
	Description: "Displays bot help",
	Run: (func(data *dcommand.Data) {
		message := &discordgo.MessageSend {
			Content: "Hi I am empty right now. Come back later",
		}
		functions.SendMessage(data.Message.ChannelID, message)
	}),
}