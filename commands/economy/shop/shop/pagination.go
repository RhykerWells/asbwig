package shop

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
	components := []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{discordgo.Button{Label: "previous", Style: 4, Disabled: true, CustomID: "shop_back"}, discordgo.Button{Label: "next", Style: 3, Disabled: true, CustomID: "shop_forward"}}}}
	guildSettings, _ := models.EconomyConfigs(qm.Where("guild_id=?", b.GuildID)).One(context.Background(), common.PQ)
	if b.MessageComponentData().CustomID != "shop_back" && b.MessageComponentData().CustomID != "shop_forward" {
		return
	}
	re := regexp.MustCompile(`\d+`)
	page, _ := strconv.Atoi(re.FindString(b.Message.Embeds[0].Footer.Text))
	if b.MessageComponentData().CustomID == "economy_forward" {
		page = page + 1
	} else {
		page = page - 1
	}
	offset :=  (page - 1) * 10
	display := ""
	guildShop, err := models.EconomyShops(qm.Where("guild_id=?", b.GuildID), qm.OrderBy("price DESC"), qm.Offset(offset)).All(context.Background(), common.PQ)
	if err != nil || len(guildShop) == 0 {
		display = "No items are in the shop for this page.\nAdd some with `createitem`"
	} else {
		display = "Buy an item with `buyitem <Name> [Quantity:Int]`\nFor more information about an item, use `iteminfo <Name>`"
		embed[0].Color = common.SuccessGreen
	}
	fields := []*discordgo.MessageEmbedField{}
	var shopNumber = 1
	for i, item := range guildShop {
		if i == 10 {
			break
		}
		shopNumber ++
		quantity := "Infinite"
		if item.Quantity.Int64 > 0 {
			quantity = humanize.Comma(item.Quantity.Int64)
		}
		price := humanize.Comma(item.Price)
		fieldName := fmt.Sprintf("%s%s - %s - %s", guildSettings.Symbol, price, item.Name, quantity)
		itemField := &discordgo.MessageEmbedField{Name: fieldName, Value: item.Description, Inline: false}
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
	if len(guildShop) > shopNumber {
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