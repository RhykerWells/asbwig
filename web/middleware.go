package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
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
		if (role.Permissions&discordgo.PermissionAdministrator != 0) || (role.Permissions&discordgo.PermissionManageServer != 0) {
            return true
        }
	}
	return false
}