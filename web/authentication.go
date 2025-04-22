package web

import (
	"encoding/json"
	"net/http"

	"github.com/RhykerWells/asbwig/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var OauthConf *oauth2.Config

func initDiscordOauth() {
	OauthConf = &oauth2.Config {
		ClientID: common.ConfigBotClientID,
		ClientSecret: common.ConfigBotSecret,
		Scopes: []string{"identify"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://discordapp.com/api/oauth2/token",
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
		},
	}
	OauthConf.RedirectURL = "https://" + common.ConfigASBWIGHost + "/confirm"
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := OauthConf.AuthCodeURL("a", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func confirmLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "a" {
		http.Redirect(w, r, "/?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	ctx := r.Context()
	token, err := OauthConf.Exchange(ctx, code)
	if err != nil {
		http.Redirect(w, r, "/?error=oauth2failure", http.StatusTemporaryRedirect)
		return
	}

	client := OauthConf.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/v10/users/@me")
	if err != nil {
		http.Redirect(w, r, "/?error=failed_to_get_user_info", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	// Decode full response into a map
	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Redirect(w, r, "/?error=failed_to_parse_user_info", http.StatusTemporaryRedirect)
		return
	}

	// Pretty print the whole map to the terminal
	userJson, _ := json.MarshalIndent(userData, "", "  ")
	logrus.Info(string(userJson))

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}