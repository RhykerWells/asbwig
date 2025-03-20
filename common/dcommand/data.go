package dcommand

import "github.com/bwmarrin/discordgo"

type Data struct {
	Session 				*discordgo.Session
	Message 				*discordgo.Message
	MessageStrippedPrefix 	*discordgo.Message
	Args    				[]string
	Handler 				*CommandHandler
}
