package events

import (
	"context"

	"github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/RhykerWells/asbwig/commands/economy"
	"github.com/RhykerWells/asbwig/common/models"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

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
	prefix.GuildPrefix(g.ID)
	economy.GuildEconomyAdd(g.ID)
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
