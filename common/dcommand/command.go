package dcommand

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type AsbwigCommand struct {
	Command			string
	Aliases			[]string
	Description 	string
	Args			[]*Args
	ArgsRequired	int
	Run				Run
	Data			*Data
}

type CommandHandler struct {
	cmdInstances 	[]AsbwigCommand
	cmdMap       	map[string]AsbwigCommand
}

type RegisteredCommand struct {
	Trigger		string
	Aliases		[]string
	Description string
	Args		[]*Args
}

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

func (c *CommandHandler) RegisteredCommands() (map[string]RegisteredCommand) {
	cmdMap := make(map[string]RegisteredCommand)
	for _, cmd := range c.cmdMap {
		rcmd := &RegisteredCommand{
			Trigger: 		 cmd.Command,
			Aliases:		 cmd.Aliases,
			Description: 	 cmd.Description,
			Args: 		 	 cmd.Args,
		}
		cmdMap[cmd.Command] = *rcmd
	}
	return cmdMap
}

type Run func(data *Data)