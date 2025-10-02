package events

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

// InitEvents registers the required event handlers to pass data to the bot
func InitEvents(s *discordgo.Session, database *sql.DB) {
	db = database

	s.AddHandler(botReady)
	s.AddHandler(guildJoin)
	s.AddHandler(guildLeave)
	s.AddHandler(messageCreate)
	s.AddHandler(guildMemberAdd)
	s.AddHandler(guildMemberLeave)

	log.Infoln("Event system initialised")
}
