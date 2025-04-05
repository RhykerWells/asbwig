package removeitem

import (
	"context"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "removeitem",
	Description: "Views the saved information about an item",
	Args: []*dcommand.Args{
		{Name: "Item", Type: dcommand.String},
		{Name: "Position", Type: dcommand.String, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Store", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		if len(data.Args) <= 0 {
			embed.Description = "No `Item` argument provided\nUse `shop [Page]` to view all items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		name := data.ArgsNotLowered[0]
		var item models.EconomyShop
		matchedItems, _ := models.EconomyShops(qm.Where("guild_id=? AND name=?", data.GuildID, name)).All(context.Background(), common.PQ)
		if len(matchedItems) == 0 {
			embed.Description = "This item doesn't exist. Use `shop` to view all items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		item = *matchedItems[0]
		if len(matchedItems) > 1 {
			if len(data.Args) <= 1 {
				embed.Description = "There are multiple items of this name. Please delete it via it's item position (use `iteminfo <Name>`)\nUse `shop [Page]` to view all items"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			index := functions.ToInt64(data.Args[1]) - 1
			if index < 0 || int(index) > len(matchedItems) {
				embed.Description = "This position doesn't exists. You can see an items item position (use `iteminfo <Name>`)\nUse `shop [Page]` to view all items"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			item = *matchedItems[index]
		}
		item.Delete(context.Background(), common.PQ)
		embed.Description = "Removed " + name + " successfully!"
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}