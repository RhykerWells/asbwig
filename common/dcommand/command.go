package dcommand

type AsbwigCommand struct {
	Command			[]string
	Description 	string
	Elevated		bool
	Run				Run
	Data			*Data
}

type CommandHandler struct {
	cmdInstances 	[]AsbwigCommand
	cmdMap       	map[string]AsbwigCommand
}

type Run func(data *Data)