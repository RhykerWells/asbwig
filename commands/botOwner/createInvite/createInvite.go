package createinvite

import (
	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
	Command:      "createinvite",
	Category:	  dcommand.CategoryOwner,
	Description:  "Creates an invite to the specified guild",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "GuildID", Type: dcommand.String},
	},
	Run: util.OwnerCommand(func(data *dcommand.Data) {
		channels, _ := common.Session.GuildChannels(data.Args[0])
		var channelID string
		for _, v := range channels {
			if v.Type == discordgo.ChannelTypeGuildText {
				channelID = v.ID
				break
			}
		}
		if channelID == "0" {
			functions.SendBasicMessage(data.ChannelID, "No available channels")
		}

		invite, _ := common.Session.ChannelInviteCreate(channelID, discordgo.Invite{
			MaxAge:    120,
			MaxUses:   1,
			Temporary: true,
			Unique:    true,
		})
		functions.SendBasicMessage(data.ChannelID, "discord.gg/"+invite.Code)
	}),
}
