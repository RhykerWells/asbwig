package iteminfo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "iteminfo",
	Category: 	 dcommand.CategoryEconomy,
	Description: "Views the saved information about an item",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
		{Name: "Position", Type: dcommand.String, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Store", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "iteminfo_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "iteminfo_forward"}}}}
		if len(data.Args) <= 0 {
			embed.Description = "No `Item` argument provided\nUse `shop [Page]` to view all items"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		page := 1
		if len(data.Args) > 1 {
			page, _ = strconv.Atoi(data.Args[1])
			if page < 1 {
				page = 1
			}
		}
		offset := (page - 1)
		name := data.ArgsNotLowered[0]
		display := ""
		matchedItems, err := models.EconomyShops(qm.Where("guild_id=? AND name=?", data.GuildID, name), qm.OrderBy("price DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
		if err != nil || len(matchedItems) == 0 {
			display = "No items are in the shop for this page.\nAdd some with `createitem`"
		} else {
			display = "If there are multiple items of the same name, use the buttons to search through them by price.\nBuy an item with `buyitem <Name> [Quantity:Int]`"
			embed.Color = common.SuccessGreen
		}
		fields := []*discordgo.MessageEmbedField{}
		for i, item := range matchedItems {
			if i == 1 {
				break
			}
			price := humanize.Comma(item.Price)
			quantity := "Infinite"
			if item.Quantity > 0 {
				quantity = humanize.Comma(item.Quantity)
			}
			role := "None"
			if item.Role != "0" {
				role = "<@&" + item.Role + ">"
			}
			field := []*discordgo.MessageEmbedField{
				{Name: "Name", Value: item.Name, Inline: true},
				{Name: "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀", Value: "⠀⠀", Inline: true},
				{Name: "Price", Value: price, Inline: true},
				{Name: "Description", Value: item.Description, Inline: true},
				{Name: "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀", Value: "⠀⠀", Inline: true},
				{Name: "Quantity", Value: quantity, Inline: true},
				{Name: "Role given", Value: role, Inline: false},
				{Name: "Reply message", Value: item.Reply, Inline: true},
				{Name: "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀", Value: "⠀⠀", Inline: true},
				{Name: "Item position", Value: fmt.Sprint(page), Inline: true},
			}
			if item.Soldby != "0" {
				soldField := &discordgo.MessageEmbedField{Name: "On market by", Value: "<@" + item.Soldby + ">"}
				field = append(field, soldField)
			}
			fields = append(fields, field...)
		}
		embed.Description = display
		embed.Fields = fields
		embed.Color = common.SuccessGreen
		embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
		if page != 1 {
			row := components[0].(discordgo.ActionsRow)
			btnPrev := row.Components[0].(discordgo.Button)
			btnPrev.Disabled = false
			row.Components[0] = btnPrev
			components[0] = row
		}
		if len(matchedItems) > page {
			row := components[0].(discordgo.ActionsRow)
			btnNext := row.Components[1].(discordgo.Button)
			btnNext.Disabled = false
			row.Components[1] = btnNext
			components[0] = row
		}
		msg, _ := common.Session.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{Embed: embed, Components: components})
		go disableButtons(msg.ChannelID, msg.ID)
	},
}

func disableButtons(channelID, messageID string) {
	time.Sleep(10 * time.Second)
	lbMessage, _ := common.Session.ChannelMessage(channelID, messageID)
	components := lbMessage.Components
	row := components[0].(*discordgo.ActionsRow)
	btnPrev := row.Components[0].(*discordgo.Button)
	btnNext := row.Components[1].(*discordgo.Button)
	btnPrev.Disabled = true
	btnNext.Disabled = true
	row.Components[0] = btnPrev
	row.Components[1] = btnNext
	components[0] = row
	message := &discordgo.MessageSend{
		Embed:      lbMessage.Embeds[0],
		Components: components,
	}
	functions.EditMessage(channelID, messageID, message)
}