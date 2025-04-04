package edititem

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)


var Command = &dcommand.AsbwigCommand{
	Command:     "edititem",
	Description: "Edits the values of an item in the shop",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
		{Name: "Option", Type: dcommand.String},
		{Name: "Value", Type: dcommand.Any},
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
		item, exists := models.EconomyShops(qm.Where("guild_id=? AND name=?", data.GuildID, name)).One(context.Background(), common.PQ)
		if exists != nil {
			embed.Description = "This item doesn't exist"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if len(data.Args) <= 1 {
			embed.Description = "No `Option` argument provided\nAvailable options are `name`, `quantity`, `price`, `description`, `reply`, `role`"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		option := data.Args[1]
		o, _ := regexp.Compile("(pric|nam|rol)e|description|(reply|quantit)y")
		if !o.MatchString(option) {
			embed.Description = "Invalid `Option` argument provided\nAvailable options are `name`, `quantity`, `price`, `description`, `reply`, `role`"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if len(data.Args) <= 2 {
			embed.Description = "No `Value` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		value := data.Args[2]
		displayValue := value
		switch option {
		case "name":
			value = data.ArgsNotLowered[2]
			itemExists, _ := models.EconomyShops(qm.Where("guild_id=? AND name=?", data.GuildID, value)).One(context.Background(), common.PQ)
			if itemExists != nil {
				embed.Description = "There is already an item with this name"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			item.Name = value
			item.Update(context.Background(), common.PQ, boil.Infer())
		case "price":
			var price int64
			if functions.ToInt64(value) <= 0 {
				embed.Description = "Invalid `Value` argument provided. Please supply a whole integer"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			price = functions.ToInt64(value)
			displayValue = humanize.Comma(price)
			item.Price = price
			item.Update(context.Background(), common.PQ, boil.Infer())
		case "quantity":
			if value != "infinite" && functions.ToInt64(value) <= 0 {
				embed.Description = "Invalid `Value` argument provided. Please supply a whole integer or `infinite` for unlimited"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			quantity := functions.ToInt64(value)
			displayValue = humanize.Comma(quantity)
			if quantity == 0 {
				displayValue = "Infinite"
			}
			item.Quantity = quantity
			item.Update(context.Background(), common.PQ, boil.Infer())
		case "description", "reply":
			value = strings.Join(data.ArgsNotLowered[1:], " ")
			for _, char := range "\"" {
				value = strings.ReplaceAll(value, string(char), "")
			}
			if utf8.RuneCountInString(value) > 200 {
				embed.Description = "Invalid `Value` argument provided. Please supply enter a string under 200 characters"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			switch option{
			case "description":
				item.Description = value
			case "reply":
				item.Reply = value
			}
			item.Update(context.Background(), common.PQ, boil.Infer())
		case "role":
			guildObj, _ := common.Session.State.Guild(data.GuildID)
			role := functions.GetRole(guildObj, value)
			if role != nil && value != "none" {
				embed.Description = "Invalid `Value` argument provided. Please supply a role ID or `none` for no role"
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			roleID := "0"
			displayValue = "None"
			if role != nil {
				roleID = role.ID
				displayValue = "<@&" + roleID + ">"
			}
			item.Role = roleID
			item.Update(context.Background(), common.PQ, boil.Infer())
		}
		embed.Color = common.SuccessGreen
		embed.Description = fmt.Sprintf("%s's `%s` has been changed to %s", name, option, displayValue)
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}