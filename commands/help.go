package commands

import (
	"fmt"

	"github.com/Ranger-4297/asbwig/bot/functions"
	"github.com/Ranger-4297/asbwig/common"
	"github.com/Ranger-4297/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var helpCmd = &dcommand.AsbwigCommand {
	Command:	[]string{"help"},
	Args:		[]*dcommand.Args{
		{Name:	"command", Type:	dcommand.String},
	},
	Description: "Displays bot help",
	Run: helpFunc,
}

func helpFunc(data *dcommand.Data) {
	command := ""
	if len(data.Args) > 0 {
		command = data.Args[0]
	}

	// Per-command help
	if command != "" {
		help(command, data.Message.ChannelID)
		return
	}

	// Generic help category
	basicEmbed := &discordgo.MessageEmbed {
		Author: &discordgo.MessageEmbedAuthor {
			Name:	fmt.Sprintf("%s help", common.Bot.Username),
			IconURL: common.Bot.AvatarURL("256"),
		},
		Color: 0x00FF7B,
	}
	message := &discordgo.MessageSend {
		Embed: basicEmbed,
	}
	functions.SendMessage(data.Message.ChannelID, message)
}

func help(command string, channelID string) {
	cmdMap := dcommand.CmdHndlr.RegisteredCommands()
	cmd, ok := cmdMap[command]
	if !ok {
		functions.SendBasicMessage(channelID, fmt.Sprintf("Command `%s` not found", command))
		return
	}
	helpEmbed := &discordgo.MessageEmbed {
		Author: &discordgo.MessageEmbedAuthor {
			Name:	fmt.Sprintf("%s help - %s", common.Bot.Username, command),
			IconURL: common.Bot.AvatarURL("256"),
		},
		Description: cmd.Description,
		Color: 0x00FF7B,
	}
	args := getArgs(cmd)
	helpEmbed.Description = cmd.Description
	if args != "" {
		helpEmbed.Description += "\n```" + args + "\n```"
	}
	message := &discordgo.MessageSend {
		Embed: helpEmbed,
	}
	functions.SendMessage(channelID, message)
}

func getArgs(command dcommand.RegisteredCommand) (str string) {
	for _, arg := range command.Args {
		str += command.Name[0]
		str += " <" + argHelp(arg) + ">\n"
	}
	return
}

func argHelp(arg *dcommand.Args) (str string) {
	argType := arg.Type.Help()
	str = fmt.Sprintf("%s:%s", arg.Name, argType)
	return
}