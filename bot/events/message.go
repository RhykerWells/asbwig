package events

import "github.com/bwmarrin/discordgo"

// messageCreate is sent when any new message is sent in the guild
// This is the primary call for context-based commands
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
}