package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"goji.io/v3/pat"
)

type CtxKey int

const (
	CtxKeyTmplData CtxKey = iota
)

// createCSRF generates a CSRF token to be used for validating requests such as logins
func createCSRF() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// setCSRF sets the csrf token in the clients web cache as a cookie
func setCSRF(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:    "asbwig_csrf",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
		Secure:  true,
	})
}

// getCSRF returns the csrf token from the clients cookies
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
		Secure:  true,
	})
	return ""
}

// setUserDataCookie sets the cookie containing the users account data
func setUserDataCookie(w http.ResponseWriter, userData map[string]interface{}) error {
	encodedValue, err := encodeUserData(userData)
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

// encodeUserData encodes the users account data into base64
func encodeUserData(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// decodeCookie decodes the base64 encoded user data into a map[string]interface{} and returns the decoded cookie
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

// checkUserCookie checks the stored browser cookie and returns the users information or an error
func checkUserCookie(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
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
		if (role.Permissions&discordgo.PermissionAdministrator != 0) || (role.Permissions&discordgo.PermissionManageServer != 0) {
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

		userData, err := checkUserCookie(w, r)
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

type TmplContextData map[string]interface{}

// baseTemplateDataMW provides the initial template data to be parsed within each page
func baseTemplateDataMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		baseData := TmplContextData{
			"HomeURL": URL,
			"Year":    time.Now().UTC().Year(),
		}
		ctx := context.WithValue(r.Context(), CtxKeyTmplData, baseData)

		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}

// userAndManagedGuildsInfoMW provides middleware to parse the current user data and the list of manageable guilds to the template data
func userAndManagedGuildsInfoMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userData, err := checkUserCookie(w, r)
		if err != nil {
			http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
			return
		}
		userID, _ := userData["id"].(string)

		guilds := getUserManagedGuilds(userID)
		guildList := make([]map[string]interface{}, 0)
		for guildID, guildName := range guilds {
			avatarURL := URL + "/static/img/icons/cross.png"
			if guild, err := common.Session.Guild(guildID); err == nil {
				if url := guild.IconURL("1024"); url != "" {
					avatarURL = url
				}
			}
			guildList = append(guildList, TmplContextData{
				"ID":     guildID,
				"Avatar": avatarURL,
				"Name":   guildName,
			})
		}

		tmplData, _ := ctx.Value(CtxKeyTmplData).(TmplContextData)
		tmplData["User"] = userData
		tmplData["ManagedGuilds"] = guildList

		ctx = context.WithValue(ctx, CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}

// currentGuildDataMW provides middleware to parse the current guilds data to the template data
func currentGuildDataMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		guildID := pat.Param(r, "server")

		retrievedGuild, _ := common.Session.Guild(guildID)

		channels, _ := common.Session.GuildChannels(guildID)
		sort.SliceStable(channels, func(i, j int) bool {
			return channels[i].Position < channels[j].Position
		})

		roles := retrievedGuild.Roles
		sort.SliceStable(roles, func(i, j int) bool {
			return roles[i].Position > roles[j].Position
		})

		member, _ := functions.GetMember(retrievedGuild.ID, common.Bot.ID)
		role := functions.HighestRole(retrievedGuild.ID, member)

		guildData := map[string]interface{}{
			"ID":                     retrievedGuild.ID,
			"Name":                   retrievedGuild.Name,
			"Avatar":                 retrievedGuild.IconURL("1024"),
			"Channels":               channels,
			"Roles":                  roles,
			"BotHighestRolePosition": role.Position,
		}
		if guildData["Avatar"] == "" {
			guildData["Avatar"] = URL + "/static/img/icons/cross.png"
		}

		tmplData, _ := ctx.Value(CtxKeyTmplData).(TmplContextData)
		tmplData["CurrentGuild"] = guildData

		ctx = context.WithValue(ctx, CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}

// urlDataMW provides middleware to parse the URL data to the template data
func urlDataMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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

		tmplData, _ := ctx.Value(CtxKeyTmplData).(TmplContextData)
		tmplData["TermsURL"] = termsURL
		tmplData["PrivacyURL"] = privacyURL

		ctx = context.WithValue(ctx, CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}

/* generic middleware setup
// genericMiddlewareNameMW provides middleware to parse XXXX data (and YYYY) to the template data
func genericMiddlewareNameMW(inner http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tmplData, _ := ctx.Value(CtxKeyTmplData).(TmplContextData)

		ctx = context.WithValue(ctx, CtxKeyTmplData, tmplData)
		inner.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(middleware)
}
*/
