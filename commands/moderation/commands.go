package moderation

import (
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/RhykerWells/durationutil"
	"github.com/RhykerWells/summit/bot/functions"
	"github.com/RhykerWells/summit/commands/util"
	"github.com/RhykerWells/summit/common"
	"github.com/RhykerWells/summit/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

// moderationBase returns the config and a bool to denote the active enabled status.
func moderationBase(guildID string) (*Config, bool) {
	config := GetConfig(guildID)

	if !config.ModerationEnabled {
		return nil, false
	}

	return config, true
}

// returns a boolean on whether the member has the current permissions to run the selected command
func hasRequiredRole(member *discordgo.Member, requiredRoles []string) bool {
	for _, role := range member.Roles {
		if slices.Contains(requiredRoles, role) {
			return true
		}
	}

	return false
}

// triggerDeletion returns the enabled status and time for the deleting the trigger
func triggerDeletion(config *Config) (bool, int64) {
	return config.ModerationTriggerDeletionEnabled, config.ModerationTriggerDeletionSeconds
}

// responseDeletion returns the enabled status and time for the deleting the response
func responseDeletion(config *Config) (bool, int64) {
	return config.ModerationResponseDeletionEnabled, config.ModerationResponseDeletionSeconds
}

// responseEmbed returns the fully-populated embed for responses
func responseEmbed(author, target *discordgo.User, action logAction) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Case type: %s", action.CaseType),
			IconURL: author.AvatarURL("1024"),
		},
		Description: fmt.Sprintf("%s has successfully %s %s :thumbsup:", author.Mention(), action.Name, target.Mention()),
		Color:       0x242429,
	}
}

var warnCommand = &dcommand.SummitCommand{
	Command:     "warn",
	Category:    dcommand.CategoryModeration,
	Aliases:     []string{""},
	Description: "Warns a user for a specified reason",
	Args: []*dcommand.Args{
		{Name: "User", Type: dcommand.User},
		{Name: "Reason", Type: dcommand.String},
	},
	ArgsRequired: 2,
	Run: func(data *dcommand.Data) {
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetMember(data.GuildID, data.Args[0])

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.MuteRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to warn this user (target has higher or equal role).")
			return
		}

		warnReason := strings.Join(data.Args[1:], " ")

		err := createCase(config, author, target, logWarn, data.ChannelID, warnReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target.User, logWarn)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var muteCommand = &dcommand.SummitCommand{
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
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetMember(data.GuildID, data.Args[0])

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.MuteRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to mute this user (target has higher or equal role).")
			return
		}

		muteRole := config.MuteRole
		_, err := functions.GetRole(config.GuildID, muteRole)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "No mute role has been setup. Please set one up on the dashboard.")
			return
		}

		muteReason := strings.Join(data.Args[2:], " ")
		duration, _ := durationutil.ToDuration(data.Args[1])
		if duration < 10*time.Minute {
			duration = 10 * time.Minute
		}

		ok = util.HasPerms(config.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionManageRoles)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `manage_roles`")
			return
		}

		err = muteUser(config, target.User.ID, duration)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Something went wrong. Is the bot role above the mute role?")
			return
		}

		err = createCase(config, author, target, logMute, data.ChannelID, muteReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			unmuteUser(config, author.User.ID, target.User.ID)
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target.User, logWarn)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var unmuteCommand = &dcommand.SummitCommand{
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
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetMember(data.GuildID, data.Args[0])

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.MuteRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to unmute this user (target has higher or equal role).")
			return
		}

		muteRole := config.MuteRole
		_, err := functions.GetRole(config.GuildID, muteRole)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "No mute role has been setup. Please set one up on the dashboard.")
			return
		}

		unmuteReason := strings.Join(data.Args[2:], " ")

		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionManageRoles)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `manage_roles`")
			return
		}

		err = unmuteUser(config, data.Author.ID, target.User.ID)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, "Something went wrong. Is the bot role above the mute role?")
			return
		}

		err = createCase(config, author, target, logUnmute, data.ChannelID, unmuteReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target.User, logUnmute)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var kickCommand = &dcommand.SummitCommand{
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
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetMember(data.GuildID, data.Args[0])

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.KickRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to kick this user (target has higher or equal role).")
			return
		}

		kickReason := strings.Join(data.Args[2:], " ")

		ok = util.HasPerms(data.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionKickMembers)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `kick_members`")
			return
		}

		err := kickUser(config, author, target, kickReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong: %s", err.Error()))
			return
		}

		err = createCase(config, author, target, logKick, data.ChannelID, kickReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target.User, logKick)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var banCommand = &dcommand.SummitCommand{
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
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetMember(data.GuildID, data.Args[0])

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.MuteRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		ok = functions.IsMemberHigher(data.GuildID, author, target)
		if !ok || target.User.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to ban this user (target has higher or equal role).")
			return
		}

		banReason := strings.Join(data.Args[2:], " ")
		duration, _ := durationutil.ToDuration(data.Args[1])
		if duration < 10*time.Minute {
			duration = 10 * time.Minute
		}

		ok = util.HasPerms(config.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionBanMembers)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `ban_members`")
			return
		}

		err := banUser(config, author, target, banReason, duration)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong: %s", err.Error()))
			return
		}

		err = createCase(config, author, target, logBan, data.ChannelID, banReason, duration)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target.User, logBan)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
var unbanCommand = &dcommand.SummitCommand{
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
		guild := functions.GetGuild(data.GuildID)
		author, _ := functions.GetMember(data.GuildID, data.Author.ID)
		target, _ := functions.GetUser(data.Args[0])
		targetMember := &discordgo.Member{
			User: target,
		}

		config, ok := moderationBase(guild.ID)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "The moderation system has not been enabled please enable it on the dashboard.")
			return
		}

		if config.ModerationLogChannel != "" {
			functions.SendBasicMessage(data.ChannelID, "Please setup a modlog channel I can access before running this command")
			return
		}

		hasRole := hasRequiredRole(author, config.BanRequiredRoles)
		if !hasRole {
			functions.SendBasicMessage(data.ChannelID, "You don't have the required roles for this command.")
			return
		}

		if target.ID == author.User.ID {
			functions.SendBasicMessage(data.ChannelID, "You don't have the correct permissions to unban this user.")
			return
		}

		unbanReason := strings.Join(data.Args[2:], " ")

		ok = util.HasPerms(config.GuildID, data.ChannelID, common.Bot.ID, discordgo.PermissionBanMembers)
		if !ok {
			functions.SendBasicMessage(data.ChannelID, "I don't have the required permissions to run this command.\nRequires: `ban_members`")
			return
		}

		err := unbanUser(config, data.Author.ID, target.ID)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong: %s", err.Error()))
			return
		}

		err = createCase(config, author, targetMember, logUnban, data.ChannelID, unbanReason)
		if err != nil {
			functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("Something went wrong creating the case: %s", err.Error()))
			return
		}

		ok, delay := triggerDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, data.Message.ID, time.Duration(delay)*time.Second)
		}

		responseEmbed := responseEmbed(author.User, target, logUnban)
		message, _ := functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: responseEmbed})
		ok, delay = responseDeletion(config)
		if ok {
			functions.DeleteMessage(data.ChannelID, message.ID, time.Duration(delay)*time.Second)
		}
	},
}
