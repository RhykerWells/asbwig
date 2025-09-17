package events

import (
	"context"

	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/models"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var scheduledJoinFunctions []func(g *discordgo.GuildCreate)

func RegisterGuildJoinfunctions(funcMap []func(g *discordgo.GuildCreate)) {
	scheduledJoinFunctions = append(scheduledJoinFunctions, funcMap...)
}

// guildJoin is called when the bot is added to a new guild (or upon connecting to the guild)
// This adds the guild to the relevant database tables
func guildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	log.WithFields(log.Fields{
		"guild":       g.ID,
		"owner":       g.OwnerID,
		"membercount": g.MemberCount,
	}).Infoln("Joined guild: ", g.Name)
	banned := isGuildBanned(g.ID)
	if banned {
		common.Session.GuildLeave(g.ID)
		return
	}

	for _, joinFunction := range scheduledJoinFunctions {
		joinFunction(g)
	}
}

var scheduledLeaveFunctions []func(g *discordgo.GuildDelete)

func RegisterGuildLeavefunctions(funcMap []func(g *discordgo.GuildDelete)) {
	scheduledLeaveFunctions = append(scheduledLeaveFunctions, funcMap...)
}

// guildLeave is called when the bot is removed from a guild
// This removes the guild from any tables that it is part of
func guildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		return
	}
	log.Infoln("Left guild: ", g.ID)

	for _, leaveFunction := range scheduledLeaveFunctions {
		leaveFunction(g)
	}
}

func isGuildBanned(guildID string) bool {
	exists, err := models.BannedGuilds(qm.Where("guild_id = ?", guildID)).Exists(context.Background(), common.PQ)
	if err != nil {
		return false
	}

	return exists
}
