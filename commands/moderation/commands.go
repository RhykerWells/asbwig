package moderation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/RhykerWells/durationutil"
	"github.com/bwmarrin/discordgo"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

// returns a boolean value on the status of the guild moderation
func isEnabled(guildID string) bool {
	config, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	return config.Enabled
}

// returns an array of required roles to run the selected command
func requireRoles(guildID, command string) []string {
	requiredRoles, _ := models.ModerationConfigRoles(qm.Where("guild_id = ?", guildID), qm.Where("action_type = ?", command)).All(context.Background(), common.PQ)
	var roleIDs []string
	for _, role := range requiredRoles {
		roleIDs = append(roleIDs, role.RoleID)
	}

	return roleIDs
}

// returns a boolean on whether the user has the current permissions to run the selected command
func hasCommandPermissions(guildID string, user *discordgo.Member, requireType string) bool {
	requiredRoles := requireRoles(guildID, requireType)
	for _, role := range user.Roles {
		if slices.Contains(requiredRoles, role) {
			return true
		}
	}
	return false
}

func triggerDeletion(guildID string) (bool, int) {
	config, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	return config.EnabledTriggerDeletion, config.SecondsToDeleteTrigger
}

func responseDeletion(guildID string) (bool, int) {
	config, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	return config.EnabledResponseDeletion, config.SecondsToDeleteResponse
}

// response returns the fully-populated embed for responses
func responseEmbed(author, target *discordgo.User, action logAction) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Case type: %s", action.CaseType),
			IconURL: author.AvatarURL("1024"),
		},
		Description: fmt.Sprintf("%s has successfully %s %s :thumbsup:", author.Mention(), action.Name, target.Mention()),
		Color: 0x242429,
	}
}

var warnCommand = &dcommand.AsbwigCommand{
	Command:     "warn",
	Category: 	 dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Warns a user for a specified reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Warn")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		target, _ := functions.GetMember(data.GuildID, data.Args[0])
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't warn yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if (!ok || target.User.ID == guild.OwnerID) {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to warn this user (target has higher role).", 10*time.Second)
			return
		}
		warnReason := strings.Join(data.Args[1:], " ")
		err := logCase(data.GuildID, author, target, logWarn, data.ChannelID, warnReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logWarn)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var muteCommand = &dcommand.AsbwigCommand{
	Command:     "mute",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Mutes a user for specified duration and reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Duration", Type: dcommand.Duration},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 3,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Mute")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		target, _ := functions.GetMember(data.GuildID, data.Args[0])
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't mute yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == guild.OwnerID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to warn this user (target has higher role).", 10*time.Second)
			return
		}
		config, _ := models.ModerationConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		muteRole := config.MuteRole.String
		if muteRole == "" {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "No mute role has been setup. Please set one up on the dashboard.", 10*time.Second)
			return
		}
		muteReason := strings.Join(data.Args[2:], " ")
		duration, _ := durationutil.ToDuration(data.Args[1])
		if duration < 10*time.Minute {
			duration = 10 * time.Minute
		}
		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionManageRoles)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `manage_roles`", 10*time.Second)
			return
		}
		_, err := getGuildModLogChannel(data.GuildID)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		err = muteUser(data.GuildID, target.User.ID, duration)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Something went wrong. Is the bot role above the mute role?", 10*time.Second)
			return
		}
		logCase(data.GuildID, author, target, logMute, data.ChannelID, muteReason, duration)
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logMute)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var unmuteCommand = &dcommand.AsbwigCommand{
	Command:     "unmute",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Unmutes a user for a specified reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Unmute")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		target, _ := functions.GetMember(data.GuildID, data.Args[0])
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't unmute yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == guild.OwnerID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to warn this user (target has higher role).", 10*time.Second)
			return
		}
		config, _ := models.ModerationConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		muteRole := config.MuteRole.String
		if muteRole == "" {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "No mute role has been setup. Please set one up on the dashboard.", 10*time.Second)
			return
		}
		unmuteReason := strings.Join(data.Args[2:], " ")
		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionManageRoles)
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `manage_roles`", 10*time.Second)
			return
		}
		_, err := getGuildModLogChannel(data.GuildID)
		if err != nil {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		err = unmuteUser(data.GuildID, data.Author.ID, target.User.ID)
		if err == errNotMember {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "This user is not in the server.", 10*time.Second)
			return
		} else if err == errNotMuted {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "This user is not muted.", 10*time.Second)
			return
		} else if err != nil {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "Something went wrong. Is the bot role above the mute role?", 10*time.Second)
			return
		}
		logCase(data.GuildID, author, target, logUnmute, data.ChannelID, unmuteReason)
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logUnmute)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var kickCommand = &dcommand.AsbwigCommand{
	Command:     "kick",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Kicks a user for a specified reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Kick")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		target, _ := functions.GetMember(data.GuildID, data.Args[0])
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't kick yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == guild.OwnerID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to warn this user (target has higher role).", 10*time.Second)
			return
		}
		kickReason := strings.Join(data.Args[1:], " ")
		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionKickMembers)
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `kick_members`", 10*time.Second)
			return
		}
		_, err := getGuildModLogChannel(data.GuildID)
		if err != nil {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		err = kickUser(data.GuildID, author.User.ID, target.User.ID, kickReason)
		if err != nil {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "This user is not in the server.", 10*time.Second)
			return
		}
		logCase(data.GuildID, author, target, logKick, data.ChannelID, kickReason)
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logKick)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var banCommand = &dcommand.AsbwigCommand{
	Command:     "ban",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Bans a user for specified duration and reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Duration", Type: dcommand.Duration},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 3,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Ban")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		targetUser, _ := functions.GetUser(data.Args[0])
		target := &discordgo.Member{
			User: targetUser,
		}
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't ban yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == guild.OwnerID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to ban this user (target has higher role).", 10*time.Second)
			return
		}
		banReason := strings.Join(data.Args[2:], " ")
		duration, _ := durationutil.ToDuration(data.Args[1])
		if duration < 10*time.Minute {
			duration = 10 * time.Minute
		}
		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionBanMembers)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `ban_members`", 10*time.Second)
			return
		}
		_, err := getGuildModLogChannel(data.GuildID)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		err = banUser(data.GuildID, author.User.ID, target.User.ID, banReason, duration)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "This user is already banned.", 10*time.Second)
			return
		}
		logCase(data.GuildID, author, target, logBan, data.ChannelID, banReason, duration)
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logBan)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var unbanCommand = &dcommand.AsbwigCommand{
	Command:     "unban",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Unbans a user for a specified reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		enabled := isEnabled(data.GuildID)
		if !enabled {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.", 10*time.Second)
			return
		}
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		ok := hasCommandPermissions(data.GuildID, author, "Unban")
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.", 10*time.Second)
			return
		}
		targetUser, _ := functions.GetUser(data.Args[0])
		target := &discordgo.Member{
			User: targetUser,
		}
		if target.User.ID == author.User.ID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You can't unban yourself.", 10*time.Second)
			return
		}
		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == guild.OwnerID {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to unban this user (target has higher role).", 10*time.Second)
			return
		}
		unbanReason := strings.Join(data.Args[1:], " ")
		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionBanMembers)
		if !ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `ban_members`", 10*time.Second)
			return
		}
		_, err := getGuildModLogChannel(data.GuildID)
		if err != nil {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command", 10*time.Second)
			return
		}
		err = unbanUser(data.GuildID, data.Author.ID, target.User.ID)
		if err == errNotBanned {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, 1*time.Second)
			functions.SendBasicMessage(data.ChannelID, "This user is not banned.", 10*time.Second)
			return
		}
		logCase(data.GuildID, author, target, logUnban, data.ChannelID, unbanReason)
		ok, delay := triggerDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}
		responseEmbed := responseEmbed(author.User, target.User, logUnban)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(data.GuildID)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}