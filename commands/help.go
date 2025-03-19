package commands

import (
	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/bwmarrin/discordgo"
)

var Help = &AsbwigCommand {
	Command:	[]string{"help"},
	Description: "Displays bot help",
	Run: (func(data *Data) {
		message := &discordgo.MessageSend {
			Content: "Hi I am empty right now. Come back later",
		}
		functions.SendMessage(data.Message.ChannelID, message)
	}),
}