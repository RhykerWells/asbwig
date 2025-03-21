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

func (c *CommandHandler) RegisterCommands(cmds ...*AsbwigCommand) {
	for _, cmd := range cmds {
		c.cmdInstances = append(c.cmdInstances, *cmd)
		for _, command := range cmd.Command {
			c.cmdMap[command] = *cmd
		}
	}
}

type Run func(data *Data)