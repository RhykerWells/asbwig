package common

import (
	"os"
	"strings"
)

var (
	ConfigBotToken    = os.Getenv("ASBWIG_TOKEN")
	ConfigBotClientID = os.Getenv("ASBWIG_CLIENTID")
	ConfigBotSecret   = os.Getenv("ASBWIG_CLIENTSECRET")

	ConfigPGHost     = os.Getenv("ASBWIG_PGHOST")
	ConfigPGDB       = os.Getenv("ASBWIG_PGDB")
	ConfigPGUsername = os.Getenv("ASBWIG_PGUSER")
	ConfigPGPassword = os.Getenv("ASBWIG_PGPASSWORD")

	ConfigASBWIGHost  = os.Getenv("ASBWIG_HOST")
	ConfigTermsURLOverride   = os.Getenv("ASBWIG_TERMSURLOVERRIDE")
	ConfigPrivacyURLOverride = os.Getenv("ASBWIG_PRIVACYURLOVERRIDE")

	ConfigBotOwner = os.Getenv("ASBWIG_OWNERID")
)

// ConfigDgoBotToken prefixes the bot token with the required "Bot " if it is not done by the host
func ConfigDgoBotToken() string {
	token := ConfigBotToken
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}
	return token
}
