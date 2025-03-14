package events

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Send message event
func botReady(s *discordgo.Session, r *discordgo.Ready) {
    guildCount := len(r.Guilds)
	log.Infof("Connected to: %d guilds", guildCount)
}