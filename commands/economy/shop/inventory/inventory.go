package inventory

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var Command = &dcommand.AsbwigCommand{
	Command:     "inventory",
	Category:    dcommand.CategoryEconomy,
	Aliases:     []string{"inv"},
	Description: "Guided create item",
	Args: []*dcommand.Args{
		{Name: "Page", Type: dcommand.Int, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username + " Inventory", IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "inventory_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "inventory_forward"}}}}
		page := 1
		if len(data.Args) > 0 {
			page, _ = strconv.Atoi(data.Args[0])
			if page < 1 {
				page = 1
			}
		}
		offset := (page - 1) * 10
		display := ""
		userInventory, err := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(data.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(data.Author.ID), qm.OrderBy("quantity DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
		if err != nil || len(userInventory) == 0 {
			display = "There are no item on this page\nBuy some with `buyitem`"
		} else {
			display = "Use an item with `useitem <Name>`\nFor more information about an item, use `iteminfo <Name>`"
			embed.Color = common.SuccessGreen
		}
		fields := []*discordgo.MessageEmbedField{}
		var invNumber = 1
		for i, item := range userInventory {
			if i == 10 {
				break
			}
			invNumber++
			role := "None"
			if item.Role != "0" {
				role = "<@&" + item.Role + ">"
			}
			itemField := &discordgo.MessageEmbedField{Name: item.Name, Value: fmt.Sprintf("Description: %s\nQuantity: %s\nRole given: %s", item.Description, humanize.Comma(item.Quantity), role), Inline: false}
			fields = append(fields, itemField)
		}
		embed.Description = display
		embed.Fields = fields
		embed.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
		if page != 1 {
			row := components[0].(discordgo.ActionsRow)
			btnPrev := row.Components[0].(discordgo.Button)
			btnPrev.Disabled = false
			row.Components[0] = btnPrev
			components[0] = row
		}
		if len(userInventory) > invNumber {
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
