package internal

import (
	"os"
	"strings"
)

var (
	ConfigBotName  = os.Getenv("asbwig_BOTNAME")
	ConfigBotToken = os.Getenv("asbwig_TOKEN")

	ConfigPGHost     = os.Getenv("asbwig_PGHOST")
	ConfigPGDB       = os.Getenv("asbwig_PGDB")
	ConfigPGUsername = os.Getenv("asbwig_PGUSER")
	ConfigPGPassword = os.Getenv("asbwig_PGPASSWORD")
)

func ConfigDgoBotToken() string {
	token := ConfigBotToken
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}
	return token
}