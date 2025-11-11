package events

import (
	"context"

	"github.com/RhykerWells/Summit/bot/core/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// scheduledGuildJoinFunctions serves as a map of all the functions that is run when a guild adds the bot
var scheduledGuildJoinFunctions []func(g *discordgo.GuildCreate)

// RegisterGuildJoinfunctions adds each guild join function to the map of functions ran when a guild adds the bot
func RegisterGuildJoinfunctions(funcMap []func(g *discordgo.GuildCreate)) {
	scheduledGuildJoinFunctions = append(scheduledGuildJoinFunctions, funcMap...)
}

// guildJoin is called when the bot is added to a new guild
// This adds the guild to the relevant database tables
func guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	banned := isGuildBanned(g.ID)
	if banned {
		s.GuildLeave(g.ID)
		return
	}

	log.WithFields(log.Fields{
		"guild":       g.ID,
		"owner":       g.OwnerID,
		"membercount": g.MemberCount,
	}).Infoln("Joined guild: ", g.Name)

	for _, joinFunction := range scheduledGuildJoinFunctions {
		go joinFunction(g)
	}
}

// scheduledGuildJoinFunctions serves as a map of all the functions that is run when a guild adds the bot
var scheduledGuildLeaveFunctions []func(g *discordgo.GuildDelete)

// RegisterGuildJoinfunctions adds each guild join function to the map of functions ran when a guild adds the bot
func RegisterGuildLeavefunctions(funcMap []func(g *discordgo.GuildDelete)) {
	scheduledGuildLeaveFunctions = append(scheduledGuildLeaveFunctions, funcMap...)
}

// guildLeave is called when the bot is removed from a guild
// This removes the guild from any tables that it is part of
func guildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return
	}

	banned := isGuildBanned(g.ID)
	if !banned {
		log.Infoln("Left guild: ", g.ID)
	}

	for _, leaveFunction := range scheduledGuildLeaveFunctions {
		go leaveFunction(g)
	}
}

// isGuildBanned returns a boolean checking against the bots banned guild database
func isGuildBanned(guildID string) bool {
	exists, err := models.BannedGuilds(models.BannedGuildWhere.GuildID.EQ(guildID)).Exists(context.Background(), db)
	if err != nil {
		return false
	}

	return exists
}
