package bot

import (
	"github.com/Ranger-4297/asbwig/bot/events"
	"github.com/bwmarrin/discordgo"
)

func Run(s *discordgo.Session) {
	events.InitEvents()
}