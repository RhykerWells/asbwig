package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"goji.io/v3/pat"
)

var sessionStore = cache.New(24*time.Hour*30, 1*time.Hour)

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
		Expires: time.Now().Add(300 * time.Second),
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
func setUserSession(w http.ResponseWriter, user *discordgo.User) {
	sessionID := uuid.NewString()
	sessionStore.Set(sessionID, user, cache.DefaultExpiration)

	http.SetCookie(w, &http.Cookie{
		Name:    "asbwig_userinfo",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24*time.Hour*30),
	})
}

// getUserSession retrieves the user data from the user session cookie
func getUserSession(sessionID string) (*discordgo.User, bool) {
    if data, found := sessionStore.Get(sessionID); found {
        return data.(*discordgo.User), true
    }
    return nil, false
}

// checkUserCookie checks the stored browser cookie and returns the users information or an error
func checkUserCookie(w http.ResponseWriter, r *http.Request) (*discordgo.User, error) {
	cookie, err := r.Cookie("asbwig_userinfo")
	if err == nil {
		// Verify cookie session
		if user, found := getUserSession(cookie.Value); found {
			return user, nil
		}

		// If verification failed — clear the bad cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "asbwig_userinfo",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
	}
	return nil, errors.New("no valid session found")
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

		user, err := checkUserCookie(w, r)
		if err != nil {
			http.Redirect(w, r, "/?error=no_access", http.StatusFound)
			return
		}

		member, err := functions.GetMember(guildIDStr, user.ID)
		if err != nil {
			http.Redirect(w, r, "/?error=no_access", http.StatusFound)
			return
		}

		managed := isUserManaged(guildIDStr, member)
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

		user, err := checkUserCookie(w, r)
		if err != nil {
			http.Redirect(w, r, "/?error=no_access", http.StatusTemporaryRedirect)
			return
		}

		guilds := getUserManagedGuilds(user.ID)
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
		tmplData["User"] = user
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
