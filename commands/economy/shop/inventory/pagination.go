package inventory

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

func Pagination(s *discordgo.Session, b *discordgo.InteractionCreate) {
	embed := []*discordgo.MessageEmbed{{Author: &discordgo.MessageEmbedAuthor{Name: b.Member.User.Username + " inventory", IconURL: b.Member.User.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}}
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "inventory_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "inventory_forward"}}}}
	guild, _ := models.EconomyConfigs(models.EconomyConfigWhere.GuildID.EQ(b.GuildID)).One(context.Background(), common.PQ)
	if b.MessageComponentData().CustomID != "inventory_back" && b.MessageComponentData().CustomID != "inventory_forward" {
		return
	}
	re := regexp.MustCompile(`\d+`)
	page, _ := strconv.Atoi(re.FindString(b.Message.Embeds[0].Footer.Text))
	if b.MessageComponentData().CustomID == "inventory_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	offset := (page - 1) * 10
	display := ""
	userInventory, err := models.EconomyUserInventories(models.EconomyUserInventoryWhere.GuildID.EQ(b.GuildID), models.EconomyUserInventoryWhere.UserID.EQ(b.Member.User.ID), qm.OrderBy("quantity DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
	if err != nil || len(userInventory) == 0 {
		display = "There are no item on this page\nBuy some with `buyitem`"
	} else {
		display = "Use an item with `useitem <Name>`\nFor more information about an item, use `iteminfo <Name>`"
		embed[0].Color = common.SuccessGreen
	}
	fields := []*discordgo.MessageEmbedField{}
	var invNumber = 1
	for i, item := range userInventory {
		if i == 10 {
			break
		}
		if i == 10 {
			break
		}
		invNumber++
		role := "None"
		if item.Role != "0" {
			role = "<@&" + item.Role + ">"
		}
		itemField := &discordgo.MessageEmbedField{Name: item.Name, Value: fmt.Sprintf("Description: %s\nQuantity: %s%s\nRole given: %s", item.Description, guild.EconomySymbol, humanize.Comma(item.Quantity), role), Inline: false}
		fields = append(fields, itemField)
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
	if len(userInventory) > invNumber {
		row := components[0].(discordgo.ActionsRow)
		btnNext := row.Components[1].(discordgo.Button)
		btnNext.Disabled = false
		row.Components[1] = btnNext
		components[0] = row
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     embed,
			Components: components,
		},
	}
	common.Session.InteractionRespond(b.Interaction, response)
}
