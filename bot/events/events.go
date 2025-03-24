package events

import (
	"github.com/RhykerWells/asbwig/common"
	log "github.com/sirupsen/logrus"
)

func InitEvents() {
	common.Session.AddHandler(botReady)
	common.Session.AddHandler(guildJoin)
	common.Session.AddHandler(guildLeave)
	common.Session.AddHandler(messageCreate)

	log.Infoln("Event system initialised")
}
