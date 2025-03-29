package dcommand

import (
	"strings"

	"github.com/RhykerWells/asbwig/bot/functions"
	prfx "github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var CmdHndlr *CommandHandler

func NewCommandHandler() *CommandHandler {
	handler := &CommandHandler{
		cmdInstances: make([]AsbwigCommand, 0),
		cmdMap:       make(map[string]AsbwigCommand),
	}
	CmdHndlr = handler
	return CmdHndlr
}

// Handles all message create events to the bot, to pass them to child functions
func (c *CommandHandler) HandleMessageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == s.State.User.ID || event.Author.Bot {
		return
	}

	prefix, ok := checkMessagePrefix(s.State.User.ID, event)
	if !ok {
		return
	}

	prefixRemoved := strings.Split(event.Content[len(prefix):], " ")
	if len(prefixRemoved) < 1 {
		return
	}

	command := strings.ToLower(prefixRemoved[0])
	args := prefixRemoved[1:]

	outer:
	for _, cmd := range c.cmdMap {
		for _, alias := range cmd.Aliases {
			if alias == command {
				command = cmd.Command
				break outer
			}
		}
	}
	cmd, ok := c.cmdMap[command]
	if !ok {
		return
	}

	data := &Data{
		Session: s,
		Args:    args,
		Handler: c,
		Message: event.Message,
	}

	go runCommand(cmd, data)
}

// Checkmessage checks a given content for the prefix or bot mention of the guild
func checkMessagePrefix(botID string, event *discordgo.MessageCreate) (prefix string, found bool) {
	prefix, ok := findBasicPrefix(event.Content, event.GuildID)
	if ok {
		return prefix, true
	}
	prefix, ok = findMentionPrefix(botID, event.Content)
	if ok {
		return prefix, true
	}
	return "", false
}

// findBasicPrefix finds a text based prefix such as "-" or "~"
func findBasicPrefix(message string, guildID string) (string, bool) {
	prefix := prfx.GuildPrefix(guildID)
	if strings.HasPrefix(message, prefix) {
		return prefix, true
	}
	return "", false
}

// findMentionPrefix finds a bot mention prefix such as @ASBWIG
func findMentionPrefix(botID string, message string) (string, bool) {
	prefix := ""
	ok := false

	if strings.Index(message, "<@"+botID+">") == 0 {
		prefix = "<@" + botID + ">"
		ok = true
	} else if strings.Index(message, "<@!"+botID+">") == 0 {
		prefix = "<@!" + botID + ">"
		ok = true
	}
	return prefix, ok
}

func runCommand(cmd AsbwigCommand, data *Data) {
	argCount := len(data.Args)
	if argCount < cmd.ArgsRequired {
		functions.SendBasicMessage(data.Message.ChannelID, "Not enough arguments passed")
		return
	}

	cmd.Run(data)

	logrus.WithFields(logrus.Fields{
		"Guild":           data.Message.GuildID,
		"Command":         cmd.Command,
		"Triggering user": data.Message.Author.ID},
	).Infoln("Executed command")
}
