package common

import (
	"os"
	"strings"
)

var (
	ConfigBotName     = os.Getenv("ASBWIG_BOTNAME")
	ConfigBotToken    = os.Getenv("ASBWIG_TOKEN")
	ConfigBotClientID = os.Getenv("ASBWIG_CLIENTID")
	ConfigBotSecret   = os.Getenv("ASBWIG_CLIENTSECRET")
	ConfigASBWIGHost  = os.Getenv("ASBWIG_HOST")

	ConfigPGHost     = os.Getenv("ASBWIG_PGHOST")
	ConfigPGDB       = os.Getenv("ASBWIG_PGDB")
	ConfigPGUsername = os.Getenv("ASBWIG_PGUSER")
	ConfigPGPassword = os.Getenv("ASBWIG_PGPASSWORD")

	ConfigBotOwner = os.Getenv("ASBWIG_OWNERID")
)

func ConfigDgoBotToken() string {
	token := ConfigBotToken
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}
	return token
}