package events

import (
	"github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/RhykerWells/asbwig/commands/economy"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// guildJoin is called when the bot is added to a new guild
// This adds the guild to the relevant database tables
func guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	prefix.GuildPrefix(g.ID)
	economy.GuildEconomyAdd(g.ID)
	log.WithFields(log.Fields{
		"guild":       g.ID,
		"owner":       g.OwnerID,
		"membercount": g.MemberCount,
	}).Infoln("Joined guild: ", g.Name)
}

// guildLeave is called when the bot is removed from a guild
// This removes the guild from any tables that it is part of
func guildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return
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
