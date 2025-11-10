package dcommand

import "github.com/bwmarrin/discordgo"

// Data defines the required data passed to each command
type Data struct {
	Session *discordgo.Session

	GuildID   string
	ChannelID string
	Author    *discordgo.User

	Message    *discordgo.Message
	ParsedArgs []*ParsedArg

	Handler *CommandHandler
}
