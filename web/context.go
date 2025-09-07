package web

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"goji.io/v3/pat"
)

// dashboardContextData returns all the necessary context data that we use within the site into one data map
func dashboardContextData(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	userData, _ := checkCookie(w, r)
	userID, _ := userData["id"].(string)

	// Set the list of guilds that the user manages
	guilds := getUserManagedGuilds(userID)
	guildList := make([]map[string]interface{}, 0)
	for guildID, guildName := range guilds {
		avatarURL := URL + "/static/img/icons/cross.png"
		if guild, err := common.Session.Guild(guildID); err == nil {
			if url := guild.IconURL("1024"); url != "" {
				avatarURL = url
			}
		}
		guildList = append(guildList, map[string]interface{}{
			"ID":   guildID,
			"Avatar": avatarURL,
			"Name": guildName,
		})
	}

	// If a guild is selected, populate a map of data
	// TODO: host goji internally and modify pat.Param to return blank string instead of panic when param is not found.
	var guildID string
	if strings.Contains(r.URL.Path, "/manage") {
		guildID = pat.Param(r, "server")
	}
	guildData := getGuildData(guildID)

	urls := getUrlData()

	// Marshal the guild data into JSON and write to the response
	responseData := map[string]interface{}{
		"User": userData,
		"ManagedGuilds": guildList,
		"Year": time.Now().UTC().Year(),
		"URLs": urls,
		"CurrentGuild": guildData,
	}
	return responseData
}

func getUrlData() (urlData map[string]interface{}) {
	u, err := url.Parse(TermsURL)
	termsURL := URL + "/terms"
	if err == nil {
		termsURL = u.String()
	}

	u, err = url.Parse(PrivacyURL)
	privacyURL := URL + "/privacy"
	if err == nil {
		privacyURL = u.String()
	}
	urlData = map[string]interface{}{
		"Home": URL,
		"Terms": termsURL,
		"Privacy": privacyURL,
	}

	return urlData
}

// getGuildData retrieves select data about the guild to use within the manage page of the dashboard
func getGuildData(guildID string) (guildData map[string]interface{}) {
	if guildID == "" {
		return guildData
	}
	retrievedGuild, _ := common.Session.Guild(guildID)
	channels, _ := common.Session.GuildChannels(guildID)
	roles := retrievedGuild.Roles
	sort.SliceStable(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})
	sort.SliceStable(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})
	moderationData := getGuildModerationSettings(guildID)
	member, _ := functions.GetMember(retrievedGuild.ID, common.Bot.ID)
	role := functions.HighestRole(retrievedGuild.ID, member)
	guildData = map[string]interface{}{
		"ID": retrievedGuild.ID,
		"Name": retrievedGuild.Name,
		"Avatar": retrievedGuild.IconURL("1024"),
		"Channels": channels,
		"Roles": roles,
		"ModerationConfig": moderationData,
		"BotHighestRolePosition": role.Position,
	}
	if guildData["Avatar"] == "" {
		guildData["Avatar"] = URL + "/static/img/icons/cross.png"
	}
	return guildData
}

func getGuildModerationSettings(guildID string) map[string]interface{} {
	config, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	commandRestrictions := getRoleRestrictions(guildID)
	triggerSettings := map[string]interface{}{
		"Enabled": config.EnabledTriggerDeletion,
		"Seconds": config.SecondsToDeleteTrigger,
	}
	responseSettings := map[string]interface{}{
		"Enabled": config.EnabledResponseDeletion,
		"Seconds": config.SecondsToDeleteResponse,
	}
	moderationSettings := map[string]interface{} {
		"Enabled": config.Enabled,
		"Modlog": config.ModLog.String,
		"Restrictions": commandRestrictions,
		"TriggerSettings": triggerSettings,
		"ResponseSettings": responseSettings,
		"MuteRole": config.MuteRole.String,
		"ManageMuteRole": config.ManageMuteRole,
		"MuteUpdateRole": config.UpdateRoles,
	}
	return moderationSettings
}

func getRoleRestrictions(guildID string) map[string][]string {
	roles, _ := models.ModerationConfigRoles(qm.Where("guild_id = ?", guildID)).All(context.Background(), common.PQ)
	commandRestrictions := make(map[string][]string)
	for _, role := range roles {
		commandRestrictions[role.ActionType] = append(commandRestrictions[role.ActionType], role.RoleID)
	}
	return commandRestrictions
}