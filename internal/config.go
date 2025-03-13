package internal

import (
    "os"
    "strings"
)

var (
    ConfigBotName           = os.Getenv("ASBWIG_BOTNAME")
    ConfigBotToken          = os.Getenv("ASBWIG_TOKEN")

    ConfigPGHost            = os.Getenv("ASBWIG_PGHOST")
    ConfigPGDB              = os.Getenv("ASBWIG_PGDB")
    ConfigPGUsername        = os.Getenv("ASBWIG_PGUSER")
    ConfigPGPassword        = os.Getenv("ASBWIG_PGPASSWORD")
)

func ConfigDgoBotToken() string {
	token := ConfigBotToken
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}
	return token
}