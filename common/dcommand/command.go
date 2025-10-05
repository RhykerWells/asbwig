package dcommand

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	CategoryGeneral = CommandCategory{
		Name:        "General",
		Description: "General bot commands",
	}
	CategoryOwner = CommandCategory{
		Name:        "Owner",
		Description: "Mainanance and other bot-owner commands",
	}
	CategoryEconomy = CommandCategory{
		Name:        "Economy",
		Description: "Gambling and other economy based commands",
	}
	CategoryModeration = CommandCategory{
		Name:        "Moderation",
		Description: "Moderation and guild safety",
	}
)

// AsbwigCommand defines the general data that must be set during the addition of a new command
type AsbwigCommand struct {
	Command      string
	Category     CommandCategory
	Aliases      []string
	Description  string
	Args         []*Args
	ArgsRequired int
	Run          Run
	Data         *Data
}

// CommandCategory defines the available category types for commands
type CommandCategory struct {
	Name        string
	Description string
}

// CommandHandler defines the general command handler, the full instances of a command and a string map to retireve them
type CommandHandler struct {
	cmdInstances []AsbwigCommand
	cmdMap       map[string]AsbwigCommand
}

// RegisteredCommand defines the context required to access data surrounding a command
type RegisteredCommand struct {
	Trigger     string
	Category    CommandCategory
	Aliases     []string
	Description string
	Args        []*Args
}

// RegisterCommands adds each command to the command handler
func (c *CommandHandler) RegisterCommands(cmds ...*AsbwigCommand) {
	for _, cmd := range cmds {
		c.cmdInstances = append(c.cmdInstances, *cmd)
		for range cmd.Command {
			if len(cmd.Aliases) > 3 {
				aliasOver := len(cmd.Aliases) - 3
				cmd.Aliases = cmd.Aliases[:len(cmd.Aliases)-aliasOver]
				logrus.Warnln(fmt.Sprintf("%s has %d too many aliases. Automatically removed the last %d.", cmd.Command, aliasOver, aliasOver))
			}
			c.cmdMap[cmd.Command] = *cmd
		}
	}
}

// RegisteredCommands returns an array of each RegisteredCommand
func (c *CommandHandler) RegisteredCommands() map[string]RegisteredCommand {
	cmdMap := make(map[string]RegisteredCommand)
	for _, cmd := range c.cmdMap {
		rcmd := &RegisteredCommand{
			Trigger:     cmd.Command,
			Category:    cmd.Category,
			Aliases:     cmd.Aliases,
			Description: cmd.Description,
			Args:        cmd.Args,
		}
		cmdMap[cmd.Command] = *rcmd
	}
	return cmdMap
}

type Run func(data *Data)
