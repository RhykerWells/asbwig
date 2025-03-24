package util

import (
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
)

func OwnerCommand(inner dcommand.Run) dcommand.Run {
	return func(data *dcommand.Data) {
		if data.Message.Author.ID == common.ConfigBotOwner {
			inner(data)
		}
	}
}
