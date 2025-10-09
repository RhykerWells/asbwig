package buyitem

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
	Command:     "buyitem",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"buy"},
	Description: "Buys an item from the shop",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
		{Name: "Position", Type: dcommand.String, Optional: true},
		{Name: "Quantity", Type: dcommand.Int, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Store", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guildConfig, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(models.EconomyUserWhere.GuildID.EQ(data.GuildID), models.EconomyUserWhere.UserID.EQ(data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		economyUserInventory, _ := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(data.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(data.Author.ID)).All(context.Background(), common.PQ)
		if len(data.Args) <= 0 {
			embed.Description = "No `Item` argument provided\nUse `shop [Page]` to view all items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		name := data.ArgsNotLowered[0]
		var item models.EconomyShop
		matchedItems, _ := models.EconomyShops(models.EconomyShopWhere.GuildID.EQ(data.GuildID), models.EconomyShopWhere.Name.EQ(name)).All(context.Background(), common.PQ)
		if len(matchedItems) == 0 {
			embed.Description = "This item doesn't exist. Use `shop` to view all items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		item = *matchedItems[0]
		if len(matchedItems) > 1 {
			if len(data.Args) <= 1 {
				embed.Description = "There are multiple items of this name. Please include it's item position after the name (use `iteminfo <Name>`)\nSyntax: `buyitem <Item> [ItemPosition] [Quantity]`\nUse `shop [Page]` to view all items"
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
		var buyQuantity int64 = 1
		argIndex := 1
		if len(matchedItems) > 1 {
			argIndex = 2
		}
		if len(data.Args) > argIndex {
			if data.Args[argIndex] == "max" || data.Args[argIndex] == "all" {
				quantity := map[string]int64{"max": (cash / item.Price), "all": item.Quantity}
				buyQuantity = quantity[data.Args[argIndex]]
				if buyQuantity == 0 {
					buyQuantity = quantity["max"]
				}
			} else if functions.ToInt64(data.Args[argIndex]) > 0 {
				buyQuantity = functions.ToInt64(data.Args[argIndex])
			}
			if item.Quantity > 0 && buyQuantity > item.Quantity {
				embed.Description = "There's not enough of this in the shop to buy that much"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
		}
		if (name == "Chicken" || name == "chicken") && buyQuantity > 1 {
			embed.Description = "You can't buy more than one chicken at a time"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		var chicken bool = false
		for _, inventoryItem := range economyUserInventory {
			if inventoryItem.Name == "Chicken" || inventoryItem.Name == "chicken" {
				chicken = true
				break
			}
		}
		if chicken {
			embed.Description = "You can't have more than one chicken in your inventory"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if item.Soldby == data.Author.ID {
			embed.Description = "You can't buy an item that you have listed"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if (item.Price * buyQuantity) > cash {
			embed.Description = "You don't have enough money to buy this item"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		newQuantity := item.Quantity - buyQuantity
		if item.Quantity > 0 && newQuantity == 0 {
			item.Delete(context.Background(), common.PQ)
		} else if item.Quantity != 0 {
			item.Quantity = newQuantity
			item.Update(context.Background(), common.PQ, boil.Infer())
		}
		newQuantity = buyQuantity
		for _, inventoryItem := range economyUserInventory {
			if inventoryItem.Name == name {
				newQuantity = inventoryItem.Quantity + buyQuantity
				break
			}
		}
		embed.Color = common.SuccessGreen
		embed.Description = fmt.Sprintf("You bought %s of %s for %s%s", humanize.Comma(buyQuantity), name, guildConfig.Symbol, humanize.Comma(item.Price*buyQuantity))
		userInventory := models.EconomyUserInventory{GuildID: data.GuildID, UserID: data.Author.ID, Name: name, Description: item.Description, Quantity: newQuantity, Role: item.Role, Reply: item.Reply}
		userInventory.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserInventoryColumns.GuildID, models.EconomyUserInventoryColumns.UserID, models.EconomyUserInventoryColumns.Name}, boil.Whitelist(models.EconomyUserInventoryColumns.Quantity), boil.Infer())
		cash = cash - (item.Price * buyQuantity)
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{models.EconomyUserColumns.GuildID, models.EconomyUserColumns.UserID}, boil.Whitelist(models.EconomyUserColumns.Cash), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}
