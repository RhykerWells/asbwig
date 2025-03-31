package events

import (
	"context"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func guildMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	guildid := m.GuildID
	userid := m.Member.User.ID
	_, _ = models.EconomyCashes(qm.Where("guild_id=?", guildid), qm.Where("user_id=?", userid)).DeleteAll(context.Background(), common.PQ)
	_, _ = models.EconomyBanks(qm.Where("guild_id=?", guildid), qm.Where("user_id=?", userid)).DeleteAll(context.Background(), common.PQ)
	_, _ = models.EconomyCooldowns(qm.Where("guild_id=?", guildid), qm.Where("user_id=?", userid)).DeleteAll(context.Background(), common.PQ)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guildid := m.GuildID
	userid := m.Member.User.ID
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", guildid)).One(context.Background(), common.PQ)
	
	cash := models.EconomyCash{
		GuildID: guildid,
		UserID: userid,
		Cash: guild.Startbalance,
	}
	_ = cash.Insert(context.Background(), common.PQ, boil.Infer())
}