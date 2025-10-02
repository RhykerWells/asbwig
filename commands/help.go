package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var helpCmd = &dcommand.AsbwigCommand{
	Command:  "help",
	Aliases:  []string{"h"},
	Category: dcommand.CategoryGeneral,
	Args: []*dcommand.Args{
		{Name: "Command", Type: dcommand.String, Optional: true},
	},
	Description: "Displays bot help",
	Run:         helpFunc,
}

// helpFunc is the entry point for the "help" command.
// If a command name is provided in Args, it will show detailed help for that command.
// Otherwise, it will show a generic category overview of all commands.
func helpFunc(data *dcommand.Data) {
	command := ""
	if len(data.Args) > 0 {
		command = data.Args[0]
	}

	// Per-command help
	if command != "" {
		help(command, data.ChannelID)
		return
	}

	// Generic help category
	genericCategoryHelp(data.ChannelID)
}

// genericCategoryHelp builds and sends an embed listing all available categories
// and their commands. The "General" category is always listed first, followed
// by the other categories sorted alphabetically.
func genericCategoryHelp(channelID string) {
	cmdMap := dcommand.CmdHndlr.RegisteredCommands()
	categories := make(map[string][]string)
	for _, cmd := range cmdMap {
		categories[cmd.Category.Name] = append(categories[cmd.Category.Name], cmd.Trigger)
	}
	categoryNames := make([]string, 0, len(categories))
	for categoryName := range categories {
		categoryNames = append(categoryNames, categoryName)
	}
	sort.SliceStable(categoryNames, func(i, j int) bool {
		if categoryNames[i] == "General" {
			return true
		}
		if categoryNames[j] == "General" {
			return false
		}
		return categoryNames[i] < categoryNames[j]
	})

	helpEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s help", common.Bot.Username),
			IconURL: common.Bot.AvatarURL("256"),
		},
		Description: "Here are the available categories and commands:",
		Color:       common.SuccessGreen,
	}
	for _, categoryName := range categoryNames {
		// Sort commands within the category
		sort.Strings(categories[categoryName])
		categoryStr := fmt.Sprintf("**%s**: `%s`", categoryName, strings.Join(categories[categoryName], "`, `"))
		helpEmbed.Description += "\n\n" + categoryStr
	}

	message := &discordgo.MessageSend{
		Embed: helpEmbed,
	}
	functions.SendMessage(channelID, message)
}

// help shows detailed help for a specific command, including its description,
// aliases, and expected arguments. If the command cannot be found, a simple
// error message is sent instead.
func help(command string, channelID string) {
	cmdMap := dcommand.CmdHndlr.RegisteredCommands()
	cmd, ok := cmdMap[command]
	if !ok {
		functions.SendBasicMessage(channelID, fmt.Sprintf("Command `%s` not found", command))
		return
	}
	helpEmbed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s help - %s/%s", common.Bot.Username, command, strings.Join(cmd.Aliases, "/")),
			IconURL: common.Bot.AvatarURL("256"),
		},
		Description: cmd.Description,
		Color:       common.SuccessGreen,
	}
	args := getArgs(cmd)
	helpEmbed.Description = cmd.Description
	if args != "" {
		helpEmbed.Description += "\n```" + cmd.Trigger + args + "\n```"
	}
	message := &discordgo.MessageSend{
		Embed: helpEmbed,
	}
	functions.SendMessage(channelID, message)
}


// getArgs builds the formatted string of arguments for a given command.
// Required arguments are enclosed in <angle brackets>, and optional arguments
// are enclosed in [square brackets].
func getArgs(command dcommand.RegisteredCommand) (str string) {
	for _, arg := range command.Args {
		if arg.Optional {
			str += " [" + argHelp(arg) + "]"
		} else {
			str += " <" + argHelp(arg) + ">"
		}
	}
	return
}


// argHelp returns a formatted string for a single argument, showing both its
// name and type.
func argHelp(arg *dcommand.Args) (str string) {
	argType := arg.Type.Help()
	str = fmt.Sprintf("%s:%s", arg.Name, argType)
	return
}
