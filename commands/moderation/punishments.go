package moderation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func RefreshMuteSettings(guildID string) {
	config, err := models.ModerationConfigs(qm.Where("guild_id = ?", guildID)).One(context.Background(), common.PQ)
	if err != nil {
		return
	}
	managed := config.ManageMuteRole
	if !managed {
		return
	}
	if config.MuteRole.String == "" {
		return
	}
	channels, _ := common.Session.GuildChannels(guildID)
	for _, channel := range channels {
		common.Session.ChannelPermissionSet(channel.ID, config.MuteRole.String,  discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
	}
}

func muteUser(guildID string, target string, duration time.Duration) error {
	moderationConfig, _ := models.ModerationConfigs(qm.Where("guild_id = ?", guildID)).One(context.Background(), common.PQ)
	err := functions.AddRole(guildID, target, moderationConfig.MuteRole.String)
	if err != nil {
		return err
	}
	rolesRemoved := []string{}
	member, _ := functions.GetMember(guildID, target)
	if len(moderationConfig.UpdateRoles) > 0 {
		roleSet := make(map[string]struct{}, len(moderationConfig.UpdateRoles))
		for _, role := range moderationConfig.UpdateRoles {
			roleSet[role] = struct{}{}
		}

		for _, userRole := range member.Roles {
			if _, exists := roleSet[userRole]; exists {
				rolesRemoved = append(rolesRemoved, userRole)
				functions.RemoveRole(guildID, target, userRole)
			}
		}
	}

	unmuteTime := time.Now().Add(duration)
	muteEntry := models.ModerationMute{
		GuildID: guildID,
		UserID: target,
		Roles: rolesRemoved,
		UnmuteAt: unmuteTime,
	}
	muteEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("unmute_at"), boil.Infer())

	scheduleUnmute(guildID, target, unmuteTime)

	return nil
}

var (
	errNotMuted = errors.New("user not muted")
	errAlreadyBanned = errors.New("user already banned")
	errNotBanned = errors.New("user not banned")
	errNotMember = errors.New("user not a member")
)


func unmuteUser(guildID string, author, target string) error {
	moderationConfig, _ := models.ModerationConfigs(qm.Where("guild_id = ?", guildID)).One(context.Background(), common.PQ)
	muteUser, err := models.ModerationMutes(qm.Where("guild_id = ?", guildID), qm.Where("user_id = ?", target)).One(context.Background(), common.PQ)
	if err != nil {
		return errNotMuted
	}

	targetMember, err := functions.GetMember(guildID, target)
	if err != nil {
		if author == common.Bot.ID {
			muteUser.Delete(context.Background(), common.PQ)
		}
		return errNotMember
	}

	for _, roleID := range muteUser.Roles {
		functions.AddRole(guildID, target, roleID)
	}
	err = functions.RemoveRole(guildID, target, moderationConfig.MuteRole.String)
	if err != nil {
		return err
	}

	muteUser.Delete(context.Background(), common.PQ)

	if author == common.Bot.ID {
		modlogChannel, _ := getGuildModLogChannel(guildID)
		botMember, _ := functions.GetMember(guildID, common.Bot.ID)
		logCase(guildID, botMember, targetMember, logUnmute, modlogChannel, "Automatic unmute")
	}

	return nil
}

func scheduleUnmute(guildID string, target string, unmuteTime time.Time) {
	delay := time.Until(unmuteTime)
	if delay <= 0 {
		go unmuteUser(guildID, common.Bot.ID, target)
		return
	}

	go func() {
		time.Sleep(time.Until(unmuteTime))
		unmuteUser(guildID, common.Bot.ID, target)
	}()
}

func scheduleAllPendingUnmutes() {
	mutes, err := models.ModerationMutes(qm.Where("unmute_at > ?", time.Now())).All(context.Background(), common.PQ)
	if err != nil {
		return
	}

	for _, mute := range mutes {
		scheduleUnmute(mute.GuildID, mute.UserID, mute.UnmuteAt)
	}
}

func kickUser(guildID, author, target, reason string) error {
	_, err := functions.GetMember(guildID, target)
	if err != nil {
		return errNotMember
	}

	authorMember, _ := functions.GetMember(guildID, author)
	auditLogReason := fmt.Sprintf("%s: %s", authorMember.User.Username, reason)

	functions.GuildKickMember(guildID, target, auditLogReason)
	return nil
}

func banUser(guildID, author, target, reason string, duration time.Duration) error {
	_, err := common.Session.GuildBan(guildID, target)
	if err == nil {
		return errAlreadyBanned
	}
	authorMember, _ := functions.GetMember(guildID, author)
	auditLogReason := fmt.Sprintf("%s: %s", authorMember.User.Username, reason)

	unbanTime := time.Now().Add(duration)
	functions.GuildBanMember(guildID, target, auditLogReason)

	banEntry := models.ModerationBan{
		GuildID: guildID,
		UserID: target,
		UnbanAt: unbanTime,
	}
	banEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("unban_at"), boil.Infer())

	scheduleUnban(guildID, target, unbanTime)
	return nil
}

func unbanUser(guildID string, author, target string) error {
	err := functions.GuildUnbanMember(guildID, target)
	if err != nil {
		return errNotBanned
	}

	targetUser, _ := functions.GetUser(target)
	targetMember := &discordgo.Member{
		User: targetUser,
	}

	if author == common.Bot.ID {
		modlogChannel, _ := getGuildModLogChannel(guildID)
		botMember, _ := functions.GetMember(guildID, common.Bot.ID)
		logCase(guildID, botMember, targetMember, logUnban, modlogChannel, "Automatic unban")
	}

	banUser, err := models.ModerationBans(qm.Where("guild_id = ?", guildID), qm.Where("user_id = ?", target)).One(context.Background(), common.PQ)

	if err != nil {
		banUser.Delete(context.Background(), common.PQ)
	}

	return nil
}

func scheduleUnban(guildID string, target string, unmuteTime time.Time) {
	delay := time.Until(unmuteTime)
	if delay <= 0 {
		go unbanUser(guildID, common.Bot.ID, target)
		return
	}

	go func() {
		time.Sleep(time.Until(unmuteTime))
		unbanUser(guildID, common.Bot.ID, target)
	}()
}