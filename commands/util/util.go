package util

import (
	"github.com/Ranger-4297/asbwig/common"
	"github.com/Ranger-4297/asbwig/common/dcommand"
)

func OwnerCommand(inner dcommand.Run) dcommand.Run {
	return func(data *dcommand.Data){
		if data.Message.Author.ID == common.ConfigBotOwner {
			inner(data)
		}
	}
}