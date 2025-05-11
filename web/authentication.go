package web

import (
	"encoding/json"
	"net/http"

	"github.com/RhykerWells/asbwig/common"
	"golang.org/x/oauth2"
)

var OauthConf *oauth2.Config

// Handles the Oauth2 scopes for sign-in and the redirect URL
func initDiscordOauth() {
	OauthConf = &oauth2.Config{
		ClientID:     common.ConfigBotClientID,
		ClientSecret: common.ConfigBotSecret,
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://discordapp.com/api/oauth2/token",
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
		},
	}
	OauthConf.RedirectURL = "https://" + common.ConfigASBWIGHost + "/confirm"
}

// handleLogin will check for the users signin cookie, if it exists, automatically log them in, if not they are redirected to the login portal
func handleLogin(w http.ResponseWriter, r *http.Request) {
	_, err := checkCookie(w, r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	csrfToken, err := createCSRF()
	if err != nil {
		return
	}
	setCSRF(w, csrfToken)
	url := OauthConf.AuthCodeURL(csrfToken, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// confirmLogin handles the successful Discord Oauth login
// and redirects users to the dashboard
func confirmLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	csrf := getCSRF(w, r)
	state := r.URL.Query().Get("state")
	if state != csrf {
		http.Redirect(w, r, "/?error=invalid_CSRF", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := OauthConf.Exchange(ctx, code)
	if err != nil {
		http.Redirect(w, r, "/?error=oauth2_failure", http.StatusTemporaryRedirect)
		return
	}

	client := OauthConf.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/v10/users/@me")
	if err != nil {
		http.Redirect(w, r, "/?error=failed_retrieving_info", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userData)

	setCookie(w, userData)

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	userCookie, err := r.Cookie("asbwig_userinfo")
	if err != nil {
		return
	}
	deleteCookie(w, userCookie)
	csrfCookie, err := r.Cookie("asbwig_csrf")
	if err != nil {
		return
	}
	deleteCookie(w, csrfCookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}