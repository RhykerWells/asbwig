package events

import (
	"github.com/Ranger-4297/asbwig/internal"
	log "github.com/sirupsen/logrus"
)

func InitEvents() {
	internal.Session.AddHandler(botReady)
	internal.Session.AddHandler(guildJoin)
	internal.Session.AddHandler(guildLeave)
	internal.Session.AddHandler(messageCreate)

	log.Infoln("Event system initialised")
}