package bot

import (
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

func Run(s *discordgo.Session) {
	events.InitEvents(s)
	commands.InitCommands(s)
	s.Identify.Intents = gatewayIntentsUsed
}