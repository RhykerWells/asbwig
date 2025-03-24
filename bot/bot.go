package bot

import (
	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/RhykerWells/asbwig/commands"
	"github.com/bwmarrin/discordgo"
)

func Run(s *discordgo.Session) {
	events.InitEvents()
	commands.InitCommands(s)
}
