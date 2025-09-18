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

var (
	errNotMuted = errors.New("user not muted")
	errNotBanned = errors.New("user not banned")
	errNotMember = errors.New("user not a member")
)

func muteUser(config *Config, targetID string, duration time.Duration) error {
	err := functions.AddRole(config.GuildID, targetID, config.MuteRole)
	if err != nil {
		return err
	}

	rolesRemoved := []string{}
	member, _ := functions.GetMember(config.GuildID, targetID)
	if len(config.MuteUpdateRoles) > 0 {
		roleSet := make(map[string]struct{}, len(config.MuteUpdateRoles))
		for _, role := range config.MuteUpdateRoles {
			roleSet[role] = struct{}{}
		}

		for _, userRole := range member.Roles {
			if _, exists := roleSet[userRole]; exists {
				rolesRemoved = append(rolesRemoved, userRole)
				functions.RemoveRole(config.GuildID, targetID, userRole)
			}
		}
	}

	unmuteTime := time.Now().Add(duration)
	muteEntry := models.ModerationMute{
		GuildID: config.GuildID,
		UserID: targetID,
		Roles: rolesRemoved,
		UnmuteAt: unmuteTime,
	}
	muteEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("unmute_at"), boil.Infer())

	scheduleUnmute(config, targetID, unmuteTime)

	return nil
}

func unmuteUser(config *Config, authorID, targetID string) error {
	mutedUser, err := models.ModerationMutes(qm.Where("guild_id = ?", config.GuildID), qm.Where("user_id = ?", targetID)).One(context.Background(), common.PQ)
	if err != nil {
		return errNotMuted
	}

	targetMember, err := functions.GetMember(config.GuildID, targetID)
	if err != nil {
		if authorID == common.Bot.ID {
			mutedUser.Delete(context.Background(), common.PQ)
		}
		return errNotMember
	}

	for _, roleID := range mutedUser.Roles {
		functions.AddRole(config.GuildID, targetID, roleID)
	}
	functions.RemoveRole(config.GuildID, targetID, config.MuteRole)

	mutedUser.Delete(context.Background(), common.PQ)

	if authorID == common.Bot.ID {
		botMember, _ := functions.GetMember(config.GuildID, common.Bot.ID)
		createCase(config, botMember, targetMember, logUnmute, config.ModerationLogChannel, "Automatic unmute")
	}

	return nil
}

func RefreshMuteSettings(config *Config) {
	if !config.MuteManageRole {
		return
	}

	if config.MuteRole == "" {
		return
	}

	channels, _ := common.Session.GuildChannels(config.GuildID)
	for _, channel := range channels {
		common.Session.ChannelPermissionSet(channel.ID, config.MuteRole, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
	}
}

func scheduleUnmute(config *Config, targetID string, unmuteTime time.Time) {
	delay := time.Until(unmuteTime)
	if delay <= 0 {
		go unmuteUser(config, common.Bot.ID, targetID)
		return
	}

	go func() {
		time.Sleep(time.Until(unmuteTime))
		unmuteUser(config, common.Bot.ID, targetID)
	}()
}

func scheduleAllPendingUnmutes() {
	mutes, err := models.ModerationMutes(qm.Where("unmute_at > ?", time.Now())).All(context.Background(), common.PQ)
	if err != nil {
		return
	}

	for _, mute := range mutes {
		config := GetConfig(mute.GuildID)
		scheduleUnmute(config, mute.UserID, mute.UnmuteAt)
	}
}

func kickUser(config *Config, author, target *discordgo.Member, reason string) error {
	auditLogReason := fmt.Sprintf("%s: %s", author.User.Username, reason)

	err := functions.GuildKickMember(config.GuildID, target.User.ID, auditLogReason)
	if err != nil {
		return err
	}

	return nil
}

func banUser(config *Config, author, target *discordgo.Member, reason string, duration time.Duration) error {
	auditLogReason := fmt.Sprintf("%s: %s", author.User.Username, reason)

	err := functions.GuildBanMember(config.GuildID, target.User.ID, auditLogReason)
	if err != nil {
		return err
	}

	unbanTime := time.Now().Add(duration)
	banEntry := models.ModerationBan{
		GuildID: config.GuildID,
		UserID: target.User.ID,
		UnbanAt: unbanTime,
	}
	banEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("unban_at"), boil.Infer())

	scheduleUnban(config, target.User.ID, unbanTime)
	return nil
}

func unbanUser(config *Config, authorID, targetID string) error {
	bannedUser, _ := models.ModerationBans(qm.Where("guild_id = ?", config.GuildID), qm.Where("user_id = ?", targetID)).One(context.Background(), common.PQ)

	err := functions.GuildUnbanMember(config.GuildID, targetID)
	if err != nil {
		return errNotBanned
	}


	targetUser, _ := functions.GetUser(targetID)
	targetMember := &discordgo.Member{
		User: targetUser,
	}

	if bannedUser != nil {
		bannedUser.Delete(context.Background(), common.PQ)
	}

	if authorID == common.Bot.ID {
		botMember, _ := functions.GetMember(config.GuildID, common.Bot.ID)
		createCase(config, botMember, targetMember, logUnban, config.ModerationLogChannel, "Automatic unban")
	}

	return nil
}

func scheduleUnban(config *Config, targetID string, unmuteTime time.Time) {
	delay := time.Until(unmuteTime)
	if delay <= 0 {
		go unbanUser(config, common.Bot.ID, targetID)
		return
	}

	go func() {
		time.Sleep(time.Until(unmuteTime))
		unbanUser(config, common.Bot.ID, targetID)
	}()
}

func scheduleAllPendingUnbans() {
	bannedUsers, err := models.ModerationBans(qm.Where("unban_at > ?", time.Now())).All(context.Background(), common.PQ)
	if err != nil {
		return
	}

	for _, bannedUser := range bannedUsers {
		config := GetConfig(bannedUser.GuildID)
		_, err := common.Session.GuildBan(bannedUser.GuildID, bannedUser.UserID)
		if err != nil {
			bannedUser.Delete(context.Background(), common.PQ)
			continue
		}

		scheduleUnban(config, bannedUser.UserID, bannedUser.UnbanAt)
	}
}