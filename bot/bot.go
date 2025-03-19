package bot

import (
	"github.com/Ranger-4297/asbwig/bot/events"
	"github.com/Ranger-4297/asbwig/commands"
	"github.com/bwmarrin/discordgo"
)

func Run(s *discordgo.Session) {
	events.InitEvents()
	commands.InitCommands(s)
}