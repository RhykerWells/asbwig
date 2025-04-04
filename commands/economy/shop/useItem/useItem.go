package useItem

import (
	"context"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)


var Command = &dcommand.AsbwigCommand{
	Command:     "useitem",
	Description: "Uses an item present in your inventory",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
		{Name: "Quantity", Type: dcommand.Int},
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
		item, exists := models.EconomyUserInventories(qm.Where("guild_id=? AND user_id=? AND name=?", data.GuildID, data.Author.ID, name)).One(context.Background(), common.PQ)
		if exists != nil {
			embed.Description = "You don't have this item\nUse `inventory [Page]` to view your items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if item.Role != "0" {
			functions.AddRole(data.GuildID, data.Author.ID, item.Role)
		}
		quantity := item.Quantity
		newQuantity := quantity - 1
		if newQuantity == 0 {
			_, _ = item.Delete(context.Background(), common.PQ)
		} else if newQuantity > 0 {
			item.Quantity = newQuantity
			_ = item.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "name"}, boil.Whitelist("quantity"), boil.Infer())
		}
		functions.SendBasicMessage(data.ChannelID, item.Reply, 10)
	},
}