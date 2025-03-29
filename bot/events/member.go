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
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guildid := m.GuildID
	userid := m.Member.User.ID
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", guildid)).One(context.Background(), common.PQ)
	balance := guild.Startbalance
	
	var cash models.EconomyCash
	cash.GuildID = guildid
	cash.UserID = userid
	cash.Cash = balance
	_ = cash.Insert(context.Background(), common.PQ, boil.Infer())
}