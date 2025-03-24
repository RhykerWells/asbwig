package ping

import (
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var Command = &dcommand.AsbwigCommand{
	Command:     []string{"ping"},
	Description: "Displays bot latency",
	Run: (func(data *dcommand.Data) {
		msg, err := functions.SendBasicMessage(data.Message.ChannelID, "Ping...")
		if err != nil {
			return
		}
		started := time.Now()
		functions.EditBasicMessage(msg.ChannelID, msg.ID, "Pong! (Edit): "+(time.Since(started)*time.Microsecond).String())
	}),
}
