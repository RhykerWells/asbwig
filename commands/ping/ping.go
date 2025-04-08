package ping

import (
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "ping",
	Category: 	 dcommand.CategoryGeneral,
	Description: "Displays bot latency",
	Run: (func(data *dcommand.Data) {
		msg, err := functions.SendBasicMessage(data.ChannelID, "Ping...")
		if err != nil {
			return
		}
		started := time.Now()
		functions.EditBasicMessage(msg.ChannelID, msg.ID, "Pong! (Edit): "+(time.Since(started)*time.Microsecond).String())
	}),
}
