package bot

import (
	"github.com/Ranger-4297/asbwig/bot/prefix"
	"github.com/Ranger-4297/asbwig/internal"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Send message event
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.ID == s.State.User.ID {
        return
    }
}

// Join guild event
func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
    prefix.GuildPrefix(g.ID)
    log.WithFields(log.Fields{
            "guild": g.ID,
            "owner": g.OwnerID,
            "membercount": g.MemberCount,
    }).Infoln("Joined guild: ", g.Name)
}

// Leave guild event
func GuildLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
    removeGuildConfig(g.ID)
    log.Infoln("Left guild: ", g.ID)
}

func removeGuildConfig(guild string) {
    const query = `
    DELETE FROM core_config WHERE guild_id=$1
    `
    internal.PQ.Exec(query, guild)
}