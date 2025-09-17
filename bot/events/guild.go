package events

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/core/models"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var scheduledGuildJoinFunctions []func(g *discordgo.GuildCreate)

func RegisterGuildJoinfunctions(funcMap []func(g *discordgo.GuildCreate)) {
	scheduledGuildJoinFunctions = append(scheduledGuildJoinFunctions, funcMap...)
}

// guildJoin is called when the bot is added to a new guild
// This adds the guild to the relevant database tables
func guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	log.WithFields(log.Fields{
		"guild":       g.ID,
		"owner":       g.OwnerID,
		"membercount": g.MemberCount,
	}).Infoln("Joined guild: ", g.Name)
	banned := isGuildBanned(g.ID)
	if banned {
		s.GuildLeave(g.ID)
		return
	}

	for _, joinFunction := range scheduledGuildJoinFunctions {
		joinFunction(g)
	}
}

var scheduledGuildLeaveFunctions []func(g *discordgo.GuildDelete)

func RegisterGuildLeavefunctions(funcMap []func(g *discordgo.GuildDelete)) {
	scheduledGuildLeaveFunctions = append(scheduledGuildLeaveFunctions, funcMap...)
}

// guildLeave is called when the bot is removed from a guild
// This removes the guild from any tables that it is part of
func guildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return
	}
	log.Infoln("Left guild: ", g.ID)

	for _, leaveFunction := range scheduledGuildLeaveFunctions {
		leaveFunction(g)
	}
}

func isGuildBanned(guildID string) bool {
	exists, err := models.BannedGuilds(qm.Where("guild_id = ?", guildID)).Exists(context.Background(), db)
	if err != nil {
		return false
	}

	return exists
}
