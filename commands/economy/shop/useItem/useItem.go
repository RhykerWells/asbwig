package useItem

import (
	"context"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "useitem",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"use"},
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
		item, exists := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(data.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(data.Author.ID), models.EconomyUserInventoryWhere.Name.EQ(name)).One(context.Background(), common.PQ)
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
			item.Delete(context.Background(), common.PQ)
		} else if newQuantity > 0 {
			item.Quantity = newQuantity
			item.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserInventoryColumns.GuildID, models.EconomyUserInventoryColumns.UserID, models.EconomyUserInventoryColumns.Name}, boil.Whitelist(models.EconomyUserInventoryColumns.Quantity), boil.Infer())
		}
		functions.SendBasicMessage(data.ChannelID, item.Reply, 10)
	},
}
