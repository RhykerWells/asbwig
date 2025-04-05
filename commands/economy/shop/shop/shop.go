package shop

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
	Command:     "shop",
	Description: "Views the shop for the server",
	Args: []*dcommand.Args{
		{Name: "Page", Type: dcommand.Int, Optional: true},
	},
	Run: func(data *dcommand.Data) {
		guild, _ := common.Session.Guild(data.GuildID)
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Shop", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "shop_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "shop_forward"}}}}
		guildSettings, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		page := 1
		if len(data.Args) > 0 {
			page, _ = strconv.Atoi(data.Args[0])
			if page < 1 {
				page = 1
			}
		}
		offset := (page - 1) * 10
		display := ""
		guildShop, err := models.EconomyShops(qm.Where("guild_id=?", data.GuildID), qm.OrderBy("price DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
		if err != nil || len(guildShop) == 0 {
			display = "No items are in the shop for this page.\nAdd some with `createitem`"
		} else {
			display = "Buy an item with `buyitem <Name> [Quantity:Int]`\nFor more information about an item, use `iteminfo <Name>`"
			embed.Color = common.SuccessGreen
		}
		fields := []*discordgo.MessageEmbedField{}
		var shopNumber = 1
		for i, item := range guildShop {
			if i == 10 {
				break
			}
			shopNumber++
			quantity := "Infinite"
			if item.Quantity > 0 {
				quantity = humanize.Comma(item.Quantity)
			}
			price := humanize.Comma(item.Price)
			fieldName := fmt.Sprintf("%s%s - %s - %s", guildSettings.Symbol, price, item.Name, quantity)
			fieldDesc := item.Description
			if item.Soldby != "0" {
				fieldDesc = fmt.Sprintf("%s\nSold by: <@%s>", item.Description, item.Soldby)
			}
			itemField := &discordgo.MessageEmbedField{Name: fieldName, Value: fieldDesc, Inline: false}
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
		if len(guildShop) > shopNumber {
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