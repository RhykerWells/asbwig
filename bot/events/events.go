package events

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func InitEvents(s *discordgo.Session) {
	s.AddHandler(botReady)
	s.AddHandler(guildJoin)
	s.AddHandler(guildLeave)
	s.AddHandler(messageCreate)
	s.AddHandler(guildMemberAdd)
	s.AddHandler(guildMemberLeave)

	log.Infoln("Event system initialised")
}
