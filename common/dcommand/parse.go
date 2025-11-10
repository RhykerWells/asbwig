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

// ParsedArg represents a single argument parsed from a command invocation.
// Each argument contains its name, expected type, and the supplied value.
type ParsedArg struct {
	Name  string
	Type  ArgumentType
	Value interface{} // raw value provided by the user
}

// String returns the argument's string representation.
// It safely converts the underlying value to a string, supporting primitive numeric types.
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

// Int64 converts and returns the argument's value as an int64.
// If the value is nil or non-numeric, it safely returns 0.
func (p *ParsedArg) Int64() int64 {
	if p.Value == nil {
		return 0
	}

	return functions.ToInt64(p.Value)
}

// User attempts to resolve the argument into a Discord user.
// If the value is nil or resolution fails, it returns nil.
func (p *ParsedArg) User() *discordgo.User {
	if p.Value == nil {
		return nil
	}

	user, _ := functions.GetUser(p.String())
	return user
}

// Member attempts to resolve the argument into a Discord guild member.
// If the value is nil or resolution fails, it returns nil.
func (p *ParsedArg) Member(guildID string) *discordgo.Member {
	if p.Value == nil {
		return nil
	}

	member, _ := functions.GetMember(guildID, p.String())
	return member
}

// BetAmount returns the argument's lowercase string value, trimmed of any surrounding whitespace.
// Typically used for betting-related arguments (e.g., "all", "half", or a specific amount).
func (p *ParsedArg) BetAmount() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

// Duration attempts to parse the argument into a time.Duration pointer.
// Returns nil if parsing fails or if the argument has no value.
func (p *ParsedArg) Duration() *time.Duration {
	if p.Value == nil {
		return nil
	}

	duration, _ := durationutil.ToDuration(p.String())
	return &duration
}

// Coin returns the argument's lowercase string value, trimmed of whitespace.
// Used for commands that require a coin flip guess, e.g., "heads" or "tails".
func (p *ParsedArg) Coin() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

// BalanceType returns the argument's lowercase string value, trimmed of whitespace.
// Used for commands that expect a balance source argument such as "bank" or "cash".
func (p *ParsedArg) BalanceType() string {
	if p.Value == nil {
		return ""
	}

	return strings.ToLower(strings.TrimSpace(p.String()))
}

// handleMissingArgs sends a message to the user notifying them that one or more required
// arguments are missing from their command invocation. Optional arguments are ignored.
func handleMissingRequiredArgs(cmd SummitCommand, data *Data) {
	args := cmd.Args
	var missingArgs []*Arg

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

// handleInvalidArgs validates each argument passed to a command against its defined type.
// If an invalid argument is found, it returns an error embed and a boolean flag set to true.
// Otherwise, it returns nil and false.
func handleInvalidArgs(cmd SummitCommand, data *Data) (*discordgo.MessageEmbed, bool) {
	for _, arg := range data.ParsedArgs {
		if !arg.Type.ValidateArg(arg, data) {
			return errorEmbed(cmd.Command, data, fmt.Sprintf("Invalid `%s` argument. Expected: `%s`", arg.Name, arg.Type.Help())), true
		}
	}

	// Return nil if no errors
	return nil, false
}

// errorEmbed constructs and returns a standardized error embed for a command execution failure.
// It includes the author's username, avatar, timestamp, and an error description.
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
