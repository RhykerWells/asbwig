package bot

import (
	"github.com/Ranger-4297/asbwig/bot/prefix"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore all messages created by the bot itself
    // This isn't required in this specific example but it's a good practice.
    if m.Author.ID == s.State.User.ID {
        return
    }
}

func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
    prefix.GuildPrefix(g.ID)
    log.WithFields(log.Fields{
            "guild": g.ID,
            "owner": g.OwnerID,
            "membercount": g.MemberCount,
    }).Infoln("Joined guild: ", g.Name)
}

