package dcommand

type AsbwigCommand struct {
	Command			[]string
	Description 	string
	ArgsRequired	int
	Run				Run
	Data			*Data
}

type CommandHandler struct {
	cmdInstances 	[]AsbwigCommand
	cmdMap       	map[string]AsbwigCommand
}

type RegisteredCommand struct {
	Name	[]string
	Description string
	Args		int
}

func (c *CommandHandler) RegisterCommands(cmds ...*AsbwigCommand) {
	for _, cmd := range cmds {
		c.cmdInstances = append(c.cmdInstances, *cmd)
		for _, command := range cmd.Command {
			c.cmdMap[command] = *cmd
		}
	}
}

func (c *CommandHandler) RegisteredCommands() (map[string]RegisteredCommand) {
	cmdMap := make(map[string]RegisteredCommand)
	for _, cmd := range c.cmdMap {
		rcmd := &RegisteredCommand{
			Name: 		 cmd.Command,
			Description: cmd.Description,
		}
		cmdMap[cmd.Command[0]] = *rcmd
	}
	return cmdMap
}

type Run func(data *Data)