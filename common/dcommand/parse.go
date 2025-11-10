package dcommand

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RhykerWells/durationutil"
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/common"
	"github.com/bwmarrin/discordgo"
)

type ParsedArg struct {
	Name  string
	Type  ArgumentType
	Value interface{}
}

func (p *ParsedArg) String() string {
	if p.Value == nil {
		return ""
	}

	switch t := p.Value.(type) {
	case string:
		return t
	case int, int32, int64, uint, uint32, uint64:
		return strconv.FormatInt(functions.ToInt64(t), 10)
	default:
		return ""
	}
}

func (p *ParsedArg) Int64() int64 {
	if p.Value == nil {
		return 0
	}

	return functions.ToInt64(p.Value)
}

func (p *ParsedArg) User() *discordgo.User {
	if p.Value == nil {
		return nil
	}

	user, _ := functions.GetUser(p.String())
	return user
}

func (p *ParsedArg) Member(guildID string) *discordgo.Member {
	if p.Value == nil {
		return nil
	}

	member, _ := functions.GetMember(guildID, p.String())
	return member
}

func (p *ParsedArg) BetAmount() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

func (p *ParsedArg) Duration() *time.Duration {
	if p.Value == nil {
		return nil
	}

	duration, _ := durationutil.ToDuration(p.String())
	return &duration
}

func (p *ParsedArg) Coin() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

func (p *ParsedArg) BalanceType() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

// handleMissingArgs sends message notifying that their are required arguments missing
func handleMissingArgs(cmd SummitCommand, data *Data) {
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

// handleInvalidArgs sends a message notifying the user that an argument they have attempted to use is invalid
func handleInvalidArgs(cmd SummitCommand, data *Data) (*discordgo.MessageEmbed, bool) {
	for _, arg := range data.ParsedArgs {
		if !arg.Type.ValidateArg(arg, data) {
			return errorEmbed(cmd.Command, data, fmt.Sprintf("Invalid `%s` argument. Expected: `%s`", cmd.Command, arg.Type.Help())), true
		}
	}

	// Return nil if no errors
	return nil, false
}

// errorEmbed returns a populated embed object to denote an error in the command execution
func errorEmbed(cmd string, data *Data, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    data.Author.Username + " - " + cmd,
			IconURL: data.Author.AvatarURL("256"),
		},
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       common.ErrorRed,
		Description: description,
	}
}
