package events

import (
	"github.com/Ranger-4297/asbwig/internal"
	log "github.com/sirupsen/logrus"
)

func InitEvents() {
	internal.Session.AddHandler(messageCreate)
	internal.Session.AddHandler(guildJoin)
	internal.Session.AddHandler(guildLeave)

	log.Infoln("Event system initialised")
}