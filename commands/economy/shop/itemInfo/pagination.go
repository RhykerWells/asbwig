package iteminfo

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Pagination(s *discordgo.Session, b *discordgo.InteractionCreate) {
	guild, _ := common.Session.Guild(b.GuildID)
	embed := []*discordgo.MessageEmbed{{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Shop", IconURL: guild.IconURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "iteminfo_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "iteminfo_forward"}}}}
	if b.MessageComponentData().CustomID != "iteminfo_back" && b.MessageComponentData().CustomID != "iteminfo_forward" {
		return
	}
	re := regexp.MustCompile(`\d+`)
	page, _ := strconv.Atoi(re.FindString(b.Message.Embeds[0].Footer.Text))
	if b.MessageComponentData().CustomID == "iteminfo_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	offset :=  (page - 1)
	display := ""
	name := b.Message.Embeds[0].Fields[0].Value
	matchedItems, err := models.EconomyShops(qm.Where("guild_id=? AND name=?", guild.ID, name), qm.OrderBy("price DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
	if err != nil || len(matchedItems) == 0 {
		display = "No items are in the shop for this page.\nAdd some with `createitem`"
	} else {
		display = "If there are multiple items of the same name, use the buttons to search through them by price.\nBuy an item with `buyitem <Name> [Quantity:Int]`"
		embed[0].Color = common.SuccessGreen
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
			{Name: "Role given", Value: role, Inline: true},
			{Name: "Reply message", Value: item.Reply, Inline: false},
		}
		if item.Soldby != "0" {
			soldField := &discordgo.MessageEmbedField{Name: "On market by", Value: "<@" + item.Soldby + ">"}
			field = append(field, soldField)
		}
		fields = append(fields, field...)
	}
	embed[0].Description = display
	embed[0].Fields = fields
	embed[0].Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page: %d", page)}
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
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
			Components: components,
		},
	}
	common.Session.InteractionRespond(b.Interaction, response)
}