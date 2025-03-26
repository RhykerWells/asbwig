package events

import (
	"context"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func guildMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	guildid := m.GuildID
	userid := m.Member.User.ID
	logrus.Warnln(guildid + " " + userid)
	_, _ = models.EconomyCashes(qm.Where("guild_id=?", guildid), qm.Where("user_id=?", userid)).DeleteAll(context.Background(), common.PQ)
}