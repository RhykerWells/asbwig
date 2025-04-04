package iteminfo

import (
	"context"
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
	Description: "Views the saved information about an item",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
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
		price := humanize.Comma(item.Price)
		quantity := "Infinite"
		if item.Quantity > 0 {
			quantity = humanize.Comma(item.Quantity)
		}
		role := "None"
		if item.Role != "0" {
			role = "<@&" + item.Role + ">"
		}
		fields := []*discordgo.MessageEmbedField{
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
			soldField := &discordgo.MessageEmbedField{Name: "On market by", Value: item.Soldby}
			fields = append(fields, soldField)
		}
		embed.Fields = fields
		embed.Color = common.SuccessGreen
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}