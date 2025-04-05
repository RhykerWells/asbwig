package dcommand

import "github.com/bwmarrin/discordgo"

type Data struct {
	Session *discordgo.Session

	GuildID   string
	ChannelID string
	Author    *discordgo.User

	Message        *discordgo.Message
	Args           []string
	ArgsNotLowered []string

	Handler *CommandHandler
}