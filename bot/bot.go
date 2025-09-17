package bot

import (
	"database/sql"

	"github.com/RhykerWells/asbwig/bot/core"
	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/RhykerWells/asbwig/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	gatewayIntentsUsed = discordgo.MakeIntent(
		discordgo.IntentGuilds |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildModeration |
		discordgo.IntentGuildVoiceStates |
		discordgo.IntentGuildPresences |
		discordgo.IntentGuildMessages |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentDirectMessages |
		discordgo.IntentDirectMessageReactions |
		discordgo.IntentMessageContent |
		discordgo.IntentGuildScheduledEvents,
	)
)


func Run(s *discordgo.Session, db *sql.DB) {
	events.InitEvents(s, db)
	core.Init()
	commands.InitCommands(s)
	s.Identify.Intents = gatewayIntentsUsed
}