package commands

import (
	"strings"

	"github.com/Ranger-4297/asbwig/bot/prefix"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func newCommandHandler() *CommandHandler {
	return &CommandHandler{
		cmdInstances: make([]AsbwigCommand, 0),
		cmdMap:       make(map[string]AsbwigCommand),
	}
}

func (c *CommandHandler) registerCommand(cmd AsbwigCommand) {
	c.cmdInstances = append(c.cmdInstances, cmd)
	for _, command := range cmd.Command {
		c.cmdMap[command] = cmd
	}
}

func (c *CommandHandler) handleMessage(s *discordgo.Session, event *discordgo.MessageCreate) {
	prefix := prefix.GuildPrefix(event.GuildID)
	if event.Author.ID == s.State.User.ID || event.Author.Bot || !strings.HasPrefix(event.Content, prefix) {
		return
	}

	split := strings.Split(event.Content[len(prefix):], " ")
	if len(split) < 1 {
		return
	}

	invoke := strings.ToLower(split[0])
	args := split[1:]

	cmd, ok := c.cmdMap[invoke]
	if !ok {
		return
	}

	data := &Data{
		Session: s,
		Args:    args,
		Handler: c,
		Message: event.Message,
	}

	go cmd.Run(data)

	logrus.WithFields(logrus.Fields{
		"Guild": event.GuildID,
		"Command": invoke,
		"Triggering user": event.Author.ID},
		).Infoln("Executed command")
}