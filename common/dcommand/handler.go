package dcommand

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	prfx "github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/RhykerWells/asbwig/common"
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
	argsNotLowered := prefixRemoved[1:]
	var args []string
	for _, arg := range argsNotLowered {
		args = append(args, strings.ToLower(arg))
	}

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
		Session:        s,
		GuildID:        event.GuildID,
		ChannelID:      event.ChannelID,
		Author:         event.Author,
		Args:           args,
		ArgsNotLowered: argsNotLowered,
		Handler:        c,
		Message:        event.Message,
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
	logrus.WithFields(logrus.Fields{
		"Guild":           data.GuildID,
		"Command":         cmd.Command,
		"Triggering user": data.Author.ID},
	).Infoln("Executed command")

	argCount := len(data.Args)
	if argCount < cmd.ArgsRequired {
		handleMissingArgs(cmd, data)
		return
	}
	if argCount > 0 {
		if embed, invalid := handleInvalidArgs(cmd, data); invalid {
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
	}

	cmd.Run(data)
}

func handleMissingArgs(cmd AsbwigCommand, data *Data) {
	args := cmd.Args
	var missingArgs []*Args
	for _, arg := range args {
		if arg.Optional {
			continue
		}
		missingArgs = append(missingArgs, arg)
	}
	var argDisplay string
	for _, arg := range missingArgs {
		argDisplay += fmt.Sprintf("<%s:%s> ", arg.Name, arg.Type.Help())
	}
	embed := errorEmbed(cmd.Command, data, fmt.Sprintf("Missing required args\n```%s %s```", cmd.Command, argDisplay))
	functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
}

func handleInvalidArgs(cmd AsbwigCommand, data *Data) (*discordgo.MessageEmbed, bool) {
	for i, arg := range cmd.Args {
		input := data.Args[i]
		errorMessage := fmt.Sprintf("Invalid `%s` arg provided.", arg.Name)
		switch reflect.TypeOf(arg.Type).String() {
		case "*dcommand.StringArg":
			return nil, false
		case "*dcommand.IntArg":
			if functions.ToInt64(input) <= 0 {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide a whole number above 0.", errorMessage)), true
			}
		case "*dcommand.UserArg":
			if _, err := functions.GetMember(data.GuildID, input); err != nil {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide a user mention or ID.", errorMessage)), true
			}
		case "*dcommand.ChannelArg":
			if _, err := functions.GetChannel(data.GuildID, input); err != nil {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide a channel mention or ID.", errorMessage)), true
			}
		case "*dcommand.BetArg":
			if functions.ToInt64(input) <= 0 && input != "all" && input != "max" {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide a whole number, `all`, or `max`.", errorMessage)), true
			}
		case "*dcommand.CoinSideArg":
			if input != "heads" && input != "tails" {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide `heads` or `tails`.", errorMessage)), true
			}
		case "*dcommand.BalanceArg":
			if input != "cash" && input != "bank" {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide `cash` or `bank`.", errorMessage)), true
			}
		case "*dcommand.ResponseType":
			if input != "work" && input != "crime" {
				return errorEmbed(cmd.Command, data, fmt.Sprintf("%s\nPlease provide `work` or `crime`", errorMessage)), true
			}
		default:
			return errorEmbed(cmd.Command, data, "Something went wrong handling the arguments."), true
		}
	}
	// Return nil if no errors
	return nil, false
}

func errorEmbed(cmd string, data *Data, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:	data.Author.Username + " - " + cmd,
			IconURL: data.Author.AvatarURL("256"),
		},
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       common.ErrorRed,
		Description: description,
	}
}