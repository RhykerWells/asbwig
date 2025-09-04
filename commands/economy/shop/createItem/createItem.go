package createitem

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var (
	activeSessions = make(map[string]string)
	activeTimers   = make(map[string]*time.Timer)
)
var Command = &dcommand.AsbwigCommand{
	Command:     "createitem",
	Category:    dcommand.CategoryEconomy,
	Description: "Guided create item",
	Args: []*dcommand.Args{
		{Name: "Name", Type: dcommand.String},
	},
	Run: util.AdminOrManageServerCommand(func(data *dcommand.Data) { itemCreation(data) }),
}

func itemCreation(data *dcommand.Data) {
	guild, _ := common.Session.Guild(data.GuildID)
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: guild.Name + " Store", IconURL: guild.IconURL("256")}, Title: "Item info", Footer: &discordgo.MessageEmbedFooter{Text: "Type cancel to cancel the setup"}, Timestamp: time.Now().Format(time.RFC3339), Color: 0x0088CC}
	_, err := models.EconomyCreateitems(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
	channel, activeSession := activeSessions[data.Author.ID]
	if err == nil || activeSession {
		functions.SendBasicMessage(data.ChannelID, fmt.Sprintf("You are already creating an item in <#%s>", channel))
		return
	}
	if len(data.Args) <= 0 || len(data.Args[0]) > 60 {
		embed.Fields = []*discordgo.MessageEmbedField{{Name: "name", Value: "⠀⠀"}}
		activeSessions[data.Author.ID] = data.ChannelID
		msg, _ := common.Session.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{Content: "Please enter a name for the item (under 60 chars)", Embed: embed})
		item := models.EconomyCreateitem{GuildID: data.GuildID, UserID: data.Author.ID, MSGID: msg.ID}
		item.Insert(context.Background(), common.PQ, boil.Infer())
		common.Session.AddHandler(handleMessageCreate)
		resetTimeout(data.GuildID, data.ChannelID, data.Author.ID)
		return
	}
	itemExists, _ := models.EconomyShops(qm.Where("guild_id=? AND name=? AND soldby=0", data.GuildID, data.ArgsNotLowered[0])).One(context.Background(), common.PQ)
	if itemExists != nil {
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Content: "Please start again and enter a name that doesn't already exist"})
		return
	}
	embed.Fields = []*discordgo.MessageEmbedField{{Name: "name", Value: data.ArgsNotLowered[0]}}
	activeSessions[data.Author.ID] = data.ChannelID
	msg, _ := common.Session.ChannelMessageSendComplex(data.ChannelID, &discordgo.MessageSend{Content: "Please enter a price for the item", Embed: embed})
	item := models.EconomyCreateitem{GuildID: data.GuildID, UserID: data.Author.ID, Name: null.StringFrom(data.ArgsNotLowered[0]), MSGID: msg.ID}
	item.Insert(context.Background(), common.PQ, boil.Infer())
	resetTimeout(data.GuildID, data.ChannelID, data.Author.ID)
	common.Session.AddHandler(handleMessageCreate)
}

func resetTimeout(guildID, channelID, userID string) {
	if timer, exists := activeTimers[userID]; exists {
		timer.Stop()
	}

	activeTimers[userID] = time.AfterFunc(2*time.Minute, func() {
		delete(activeSessions, userID)
		models.EconomyCreateitems(qm.Where("guild_id=? AND user_id=?", guildID, userID)).DeleteAll(context.Background(), common.PQ)
		functions.SendBasicMessage(channelID, "The item creation session has timed out due to inactivity. Please try again")
	})
}

func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", m.GuildID)).One(context.Background(), common.PQ)
	channelID, exists := activeSessions[m.Author.ID]
	if !exists || m.ChannelID != channelID {
		return
	}
	if m.Content == "cancel" {
		delete(activeSessions, m.Author.ID)
		models.EconomyCreateitems(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.Author.ID)).DeleteAll(context.Background(), common.PQ)
		functions.SendBasicMessage(m.ChannelID, "Create item cancelled")
		return
	}
	resetTimeout(m.GuildID, m.ChannelID, m.Author.ID)
	delay := 10 * time.Second
	createItem, _ := models.EconomyCreateitems(qm.Where("guild_id=? AND user_id=?", m.GuildID, m.Author.ID)).One(context.Background(), common.PQ)
	message, _ := common.Session.ChannelMessage(m.ChannelID, createItem.MSGID)
	embed := message.Embeds[0]
	if !createItem.Name.Valid {
		name := strings.Split(m.Content, " ")[0]
		if len(name) > 60 {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "Please enter a name for the item (under 60 chars)"}, delay)
			return
		}
		itemExists, _ := models.EconomyShops(qm.Where("guild_id=? AND name=? AND soldby=0", m.GuildID, name)).One(context.Background(), common.PQ)
		if itemExists != nil {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "Please enter a name that doesn't already exist"}, delay)
			return
		}
		createItem.Name = null.StringFrom(name)
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("name"))
		embed.Fields = []*discordgo.MessageEmbedField{{Name: "name", Value: name, Inline: true}}
		functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "Please enter a price for this item", Embed: embed})
		return
	}
	if !createItem.Price.Valid {
		price := strings.Split(m.Content, " ")[0]
		if functions.ToInt64(price) <= 0 {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "Please enter a price for this item"}, delay)
			return
		}
		createItem.Price = null.Int64From(functions.ToInt64(price))
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("price"))
		priceField := &discordgo.MessageEmbedField{Name: "price", Value: fmt.Sprintf("%s%s", guild.Symbol, humanize.Comma(functions.ToInt64(price))), Inline: true}
		embed.Fields = append(embed.Fields, priceField)
		functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "Please enter a description for this item (under 200 chars)", Embed: embed})
		return
	}
	if !createItem.Description.Valid {
		if utf8.RuneCountInString(m.Content) > 200 {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "Please enter a description for this item (under 200 chars)"}, delay)
			return
		}
		createItem.Description = null.StringFrom(m.Content)
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("description"))
		descriptionField := &discordgo.MessageEmbedField{Name: "Description", Value: m.Content}
		embed.Fields = append(embed.Fields, descriptionField)
		functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "How much of this item should the store stock?\nType `skip` or `inf` to skip this step", Embed: embed})
		return
	}
	if !createItem.Quantity.Valid {
		quantity := strings.Split(m.Content, " ")[0]
		if quantity != "skip" && quantity != "inf" && functions.ToInt64(quantity) <= 0 {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "How much of this item should the store stock?\nType `skip` or `inf` to skip this step"}, delay)
			return
		}
		displayQuantity := quantity
		if quantity == "skip" || quantity == "inf" {
			displayQuantity = "Infinite"
		}
		createItem.Quantity = null.Int64From(functions.ToInt64(quantity))
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("quantity"))
		quantityField := &discordgo.MessageEmbedField{Name: "Stock", Value: displayQuantity}
		embed.Fields = append(embed.Fields, quantityField)
		functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "What role should be given when this item is used? (Role ID/Mention)\nType `skip` to skip this step", Embed: embed})
		return
	}
	if !createItem.Role.Valid {
		roleID := "0"
		displayRole := ""
		role, _ := functions.GetRole(m.GuildID, strings.Split(m.Content, " ")[0])
		if strings.Split(m.Content, " ")[0] != "skip" && role == nil {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "What role should be given when this item is used? (Role ID/Mention)\nType `skip` to skip this step"}, delay)
			return
		}
		if role != nil {
			roleID = role.ID
			displayRole = fmt.Sprintf("<@&%s>", roleID)
		}
		createItem.Role = null.StringFrom(roleID)
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("role"))
		roleField := &discordgo.MessageEmbedField{Name: "Role given", Value: displayRole}
		embed.Fields = append(embed.Fields, roleField)
		functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "What reply should be given when this item is used?", Embed: embed})
		return
	}
	if !createItem.Reply.Valid {
		if utf8.RuneCountInString(m.Content) > 200 {
			functions.DeleteMessage(m.ChannelID, m.ID)
			functions.SendMessage(m.ChannelID, &discordgo.MessageSend{Content: "Please enter a reply message for when this item is used (under 200 chars)"}, delay)
			return
		}
		createItem.Reply = null.StringFrom(m.Content)
		createItem.Update(context.Background(), common.PQ, boil.Whitelist("reply"))
		replyField := &discordgo.MessageEmbedField{Name: "Reply message", Value: m.Content}
		embed.Fields = append(embed.Fields, replyField)
	}
	embed.Footer = nil
	embed.Color = common.SuccessGreen
	item := models.EconomyShop{
		GuildID:     createItem.GuildID,
		Name:        createItem.Name.String,
		Description: createItem.Description.String,
		Price:       createItem.Price.Int64,
		Quantity:    createItem.Quantity.Int64,
		Role:        createItem.Role.String,
		Reply:       createItem.Reply.String,
		Soldby:      "0",
	}
	item.Insert(context.Background(), common.PQ, boil.Infer())
	delete(activeSessions, m.Author.ID)
	if timer, exists := activeTimers[m.Author.ID]; exists {
		timer.Stop()
	}
	models.EconomyCreateitems(qm.Where("guild_id=? AND user=?", m.GuildID, m.Author.ID)).DeleteAll(context.Background(), common.PQ)
	functions.EditMessage(m.ChannelID, createItem.MSGID, &discordgo.MessageSend{Content: "Item created! ✅", Embed: embed})
}
