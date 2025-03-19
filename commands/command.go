package commands

import "github.com/bwmarrin/discordgo"

type AsbwigCommand struct {
	Command			[]string
	Description 	string
	Elevated		bool
	Run				Run
	Data			*Data
}

type Data struct {
	Session *discordgo.Session
	Message *discordgo.Message
	Args    []string
	Handler *CommandHandler
}

type CommandHandler struct {
	cmdInstances []AsbwigCommand
	cmdMap       map[string]AsbwigCommand
}

type Run func(data *Data)