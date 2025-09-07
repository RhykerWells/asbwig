package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/bot/prefix"
	"github.com/RhykerWells/asbwig/commands/moderation"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"goji.io/v3/pat"
)

func createCSRF() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func setCSRF(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "asbwig_csrf",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
		Secure: true,
	})
}

func getCSRF(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("asbwig_csrf")
	if err == nil {
		return cookie.Value
	}

	// If decoding failed — clear the bad cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "asbwig_csrf",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		Secure: true,
	})
	return ""
}

// setCookie sets the cookie containing the users Oauth2 scope information
func setCookie(w http.ResponseWriter, userData map[string]interface{}) error {
	encodedValue, err := encodeCookie(userData)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "asbwig_userinfo",
		Value:   encodedValue,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})

	return nil
}

func encodeCookie(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func decodeCookie(encoded string) (map[string]interface{}, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(decodedBytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// checkCookie checks the stored browser cookie and returns the users information or an error
func checkCookie(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	cookie, err := r.Cookie("asbwig_userinfo")
	if err == nil {
		// Decode and verify cookie
		userData, err := decodeCookie(cookie.Value)
		if err == nil {
			return userData, nil
		}

		// If decoding failed — clear the bad cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "asbwig_userinfo",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
	}
	return nil, errors.New("no session ")
}

// deleteCookie deletes the specified HTTP cookie from local storage
func deleteCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Value = "none"
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

// getUserManagedGuilds returns the guild IDs of the guilds that the bot is in
// where the user has Owner, Manage Server or Administrator permissions
func getUserManagedGuilds(userID string) map[string]string {
	managedGuilds := make(map[string]string)
	for _, guild := range common.Session.State.Guilds {
		member, err := common.Session.GuildMember(guild.ID, userID)
		if err != nil {
			continue
		}
		managed := isUserManaged(guild.ID, member)
		if managed {
			// Store the guild ID and name in the map
			managedGuilds[guild.ID] = guild.Name
		}
	}

	return managedGuilds
}

// isUserManaged returns a boolean of whether or not the user has the permissions to manage the guild
// Permissions required are: Owner, Manage Server or Administrator 
func isUserManaged(guildID string, member *discordgo.Member) bool {
	guild, err := common.Session.State.Guild(guildID)
	if err == nil && guild.OwnerID == member.User.ID {
		return true
	}
	for _, roleID := range member.Roles {
		role, err := common.Session.State.Role(guildID, roleID)
		if err == nil {
			continue
		}
		if (role.Permissions&discordgo.PermissionAdministrator != 0) || (role.Permissions&discordgo.PermissionManageGuild != 0) {
			return true
		}
	}
	return false
}

// validateGuild ensures users can't access the manage page for guilds without the correct permissions
func validateGuild(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		guildIDStr := pat.Param(r, "server")
		_, err := strconv.ParseInt(guildIDStr, 10, 64)
		if err != nil {
			http.Redirect(w, r, "/?error=invalid_guild", http.StatusFound)
			return
		}

		userData, err := checkCookie(w, r)
		if err != nil {
			http.Redirect(w, r, "/?error=no_access", http.StatusFound)
			return
		}

		userID, _ := userData["id"].(string)
		user, err := functions.GetMember(guildIDStr, userID)
		if err != nil {
			http.Redirect(w, r, "/?error=no_access", http.StatusFound)
			return
		}

		managed := isUserManaged(guildIDStr, user)
		if !managed {
			http.Redirect(w, r, "/?error=no_access", http.StatusFound)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

// handleUpdatePrefix changes the guilds prefix in the database with the one provided from the dashboard
func handleUpdatePrefix(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data struct {
		Prefix string `json:"prefix"`
	}
	
	json.NewDecoder(r.Body).Decode(&data)

	server := pat.Param(r, "server")

	prefix.ChangeGuildPrefix(server, data.Prefix)

	http.Redirect(w, r, "/dashboard/"+server+"/manage/core", http.StatusSeeOther)
}

func handleUpdateModeration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data struct {
		Update string `json:"update"`
		Enabled bool `json:"enabled"`
		Modlog string `json:"modlog"`
		TriggerStatus bool `json:"triggerStatus"`
		TriggerInput int `json:"triggerInput"`
		ResponseStatus bool `json:"responseStatus"`
		ResponseInput int `json:"responseInput"`
		Roles map[string][]string `json:"roles"`
		MuteRole string `json:"muteRole"`
		ManagedMuteRole bool `json:"managedMuteRole"`
		UpdateRoles []string `json:"updateRoles"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	server := pat.Param(r, "server")

	config, _ := models.ModerationConfigs(qm.Where("guild_id=?", server)).One(context.Background(), common.PQ)
	switch data.Update {
	case "all":
		config.ModLog = null.StringFrom(data.Modlog)
		config.UpdateRoles = data.UpdateRoles
		config.MuteRole = null.NewString(data.MuteRole, true)
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("mod_log", "update_roles", "mute_role"), boil.Infer())
		updateAllRoles(server, data.Roles)
	case "modlog":
		config.ModLog = null.StringFrom(data.Modlog)
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("mod_log"), boil.Infer())
	case "Warn", "Mute", "Unmute", "Kick", "Ban", "Unban":
		whitelist := "required_" + strings.ToLower(data.Update) +"_roles"
		actionRoles := data.Roles[data.Update]
		updateRoles(server, data.Update, actionRoles)
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist(whitelist), boil.Infer())
	case "status":
		config.Enabled = data.Enabled
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("enabled"), boil.Infer())
	case "triggerStatus", "responseStatus":
		whitelist := ""
		switch data.Update {
		case "triggerStatus":
			whitelist = "enabled_trigger_deletion"
			config.EnabledTriggerDeletion = data.TriggerStatus
		case "responseStatus":
			whitelist = "enabled_response_deletion"
			config.EnabledResponseDeletion = data.ResponseStatus
		}
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist(whitelist), boil.Infer())
	case "triggerInput", "responseInput", "manageMuteRole":
		whitelist := ""
		switch data.Update {
		case "triggerInput":
			whitelist = "seconds_to_delete_trigger"
			config.SecondsToDeleteTrigger = data.TriggerInput
		case "responseInput":
			whitelist = "seconds_to_delete_response"
			config.SecondsToDeleteResponse = data.ResponseInput
		case "manageMuteRole":
			whitelist = "manage_mute_role"
			config.ManageMuteRole = data.ManagedMuteRole
		}
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist(whitelist), boil.Infer())
	case "updateRoles":
		config.UpdateRoles = data.UpdateRoles
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("update_roles"), boil.Infer())
	case "muteRole":
		config.MuteRole = null.NewString(data.MuteRole, true)
		config.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("mute_role"), boil.Infer())
	}

	moderation.RefreshMuteSettings(server)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}


func updateAllRoles(guildID string, rolesMap map[string][]string) {
	tx, _ := common.PQ.BeginTx(context.Background(), nil)
	for actionType, roleIDs := range rolesMap {
		_, err := models.ModerationConfigRoles(qm.Where("guild_id = ?", guildID), qm.Where("action_type = ?", actionType)).DeleteAll(context.Background(), tx)

		if err != nil {
			tx.Rollback()
			break
		}

		for _, roleID := range roleIDs {
			role := models.ModerationConfigRole{
				GuildID:   guildID,
				ActionType: actionType,
				RoleID:    roleID,
			}
			role.Insert(context.Background(), tx, boil.Infer())
		}
	}

	tx.Commit()
}

func updateRoles(guildID string, roleType string, rolesMap []string) {
	tx, _ := common.PQ.BeginTx(context.Background(), nil)
	_, err := models.ModerationConfigRoles(qm.Where("guild_id = ?", guildID), qm.Where("action_type = ?", roleType)).DeleteAll(context.Background(), tx)

	if err != nil {
		tx.Rollback()
		return
	}

	for _, roleID := range rolesMap {
		role := models.ModerationConfigRole{
			GuildID:   guildID,
			ActionType: roleType,
			RoleID:    roleID,
		}
		role.Insert(context.Background(), tx, boil.Infer())
	}

	tx.Commit()
}

func handleModerationAction(w http.Response, r *http.Request) {
	
}