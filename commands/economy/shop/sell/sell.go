package sell

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "sell",
	Category:    dcommand.CategoryEconomy,
	Description: "Adds an item from your inventory to the shop for others to buy",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
		{Name: "Price", Type: dcommand.Int},
		{Name: "Quantity", Type: dcommand.Int, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Store", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guildConfig, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		if len(data.Args) <= 0 {
			embed.Description = "No `Item` argument provided\nUse `inventory [Page]` to view all your items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		name := data.ArgsNotLowered[0]
		inventoryItem, exists := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(data.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(data.Author.ID), models.EconomyUserInventoryWhere.Name.EQ(name)).One(context.Background(), common.PQ)
		if exists != nil {
			embed.Description = "You don't have this item\nUse `inventory [Page]` to view all your items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if inventoryItem.Name == "Chicken" || inventoryItem.Name == "chicken" {
			embed.Description = "You can't sell this item\nUse `inventory [Page]` to view all your items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if len(data.Args) <= 1 {
			embed.Description = "No `Price` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		price := functions.ToInt64(data.Args[1])
		if price <= 0 {
			embed.Description = "Invalid `Price` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		sellQuantity := int64(1)
		if len(data.Args) > 2 {
			sellQuantity = functions.ToInt64(data.Args[2])
			if data.Args[2] != "all" && sellQuantity <= 0 {
				embed.Description = "Invalid `Quantity` argument provided. Please supply a whole integer or `all` for all of them"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			if data.Args[2] == "all" {
				sellQuantity = inventoryItem.Quantity
			}
		}
		item := models.EconomyShop{GuildID: data.GuildID, Name: inventoryItem.Name, Description: inventoryItem.Description, Price: price, Quantity: sellQuantity, Role: inventoryItem.Role, Reply: inventoryItem.Reply, Soldby: data.Author.ID}
		item.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "name", "soldby"}, boil.Whitelist("description", "price", "quantity", "role", "reply"), boil.Infer())
		quantity := inventoryItem.Quantity
		newQuantity := quantity - 1
		if newQuantity == 0 {
			inventoryItem.Delete(context.Background(), common.PQ)
		} else if newQuantity > 0 {
			inventoryItem.Quantity = newQuantity
			inventoryItem.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "name"}, boil.Whitelist("quantity"), boil.Infer())
		}
		embed.Description = fmt.Sprintf("Added %s to the shop. Selling for %s%s!", name, guildConfig.Symbol, humanize.Comma(price))
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
