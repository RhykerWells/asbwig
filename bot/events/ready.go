package events

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// botReady is sent when the bot originally connects to the gateway
// This is used to test the bot actually connects
func botReady(s *discordgo.Session, r *discordgo.Ready) {
	guildCount := len(r.Guilds)
	log.Infof("Connected to: %d guilds", guildCount)
}