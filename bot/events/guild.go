package events

import (
	"github.com/Ranger-4297/asbwig/bot/prefix"
	"github.com/Ranger-4297/asbwig/common"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Join guild event
func guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	prefix.GuildPrefix(g.ID)
	log.WithFields(log.Fields{
		"guild":       g.ID,
		"owner":       g.OwnerID,
		"membercount": g.MemberCount,
	}).Infoln("Joined guild: ", g.Name)
}

// Leave guild event
func guildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return // Guild outage
	}
	removeGuildConfig(g.ID)
	log.Infoln("Left guild: ", g.ID)
}

func removeGuildConfig(guild string) {
	const query = `
    DELETE FROM core_config WHERE guild_id=$1
    `
	common.PQ.Exec(query, guild)
}
