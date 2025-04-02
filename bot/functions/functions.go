package functions

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
)

// Message functions
func SendBasicMessage(channelID string, message string) (msg *discordgo.Message, err error) {
	msg, err = common.Session.ChannelMessageSend(channelID, message)
	return msg, err
}

func SendMessage(channelID string, messageData *discordgo.MessageSend, delay ...any) error {
	message, err := common.Session.ChannelMessageSendComplex(channelID, messageData)
	if delay != nil {
		DeleteMessage(channelID, message.ID, delay)
	}
	return err
}

func SendDM(userID string, messageData *discordgo.MessageSend) error {
	channel, err := common.Session.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	err = SendMessage(channel.ID, messageData)
	return err
}

func EditBasicMessage(channelID, messageID, message string) error {
	_, err := common.Session.ChannelMessageEdit(channelID, messageID, message)
	return err
}

func EditMessage(channelID string, messageID string, message *discordgo.MessageSend) {
	edit := &discordgo.MessageEdit{
		ID:	messageID,
		Channel: channelID,
	}
	if message.Content != "" {edit.Content = &message.Content}
	if message.Embed != nil {edit.Embed = message.Embed}
	if message.Embeds != nil {edit.Embeds = &message.Embeds}
	if message.Components != nil {edit.Components = &message.Components}
	_, _ = common.Session.ChannelMessageEditComplex(edit)
}

func DeleteMessage(channelID, messageData string, delay ...any) error {
	var duration int
	if len(delay) > 0 {
		duration = delay[0].([]any)[0].(int)
    }
	time.Sleep(time.Duration(duration * int(time.Second)))
	err := common.Session.ChannelMessageDelete(channelID, messageData)
	return err
}

// User functions
func GetUser(user string) (*discordgo.User, error) {
	u, err := common.Session.User(user)

	return u, err
}

func GetMember(guild string, user string) (*discordgo.Member, error) {
	// Direct mention
	if strings.HasPrefix(user, "<@") {
		user = user[2 : len(user)-1]
	}
	u, err := common.Session.GuildMember(guild, user)

	return u, err
}

// Role functions
func AddRole(guild *discordgo.Guild, member *discordgo.Member, roleID string) error {
	for _, v := range member.Roles {
		if v == roleID {
			// Already has the role
			return nil
		}
	}

	return common.Session.GuildMemberRoleAdd(guild.ID, member.User.ID, roleID)
}

func RemoveRole(guild *discordgo.Guild, member *discordgo.Member, roleID string) error {
	for _, v := range member.Roles {

		if GetRole(guild, v).ID != roleID {
			common.Session.GuildMemberRoleRemove(guild.ID, member.User.ID, roleID)
			return nil
		}
	}
	return nil
}

func SetRoles(guild *discordgo.Guild, member *discordgo.Member, roleIDs []string) error {
	roles := make(map[string]struct{})

	for _, id := range member.Roles {
		role := GetRole(guild, id)
		if role != nil && role.Managed {
			roles[id] = struct{}{}
		}
	}
	roleSlice := make([]string, 0, len(roles))
	for id := range roles {
		roleSlice = append(roleSlice, id)
	}
	userData := &discordgo.GuildMemberParams{
		Roles: &roleSlice,
	}
	_, err := common.Session.GuildMemberEdit(guild.ID, member.User.ID, userData)
	return err
}

// Misc
func SetStatus(statusText string) {
	// TODO VERSION on nothing
	if statusText == "" {
		statusText = ""
	}

	common.Session.UpdateCustomStatus(statusText)
}

// Helper tools

// ToInt64 takes the value of an int, float or string and returns it as a whole 64-bit integer if possible.
func ToInt64(conv any) int64 {
	t := reflect.ValueOf(conv)
	switch {
	case t.CanInt():
		return t.Int()
	case t.CanFloat():
		if t.Float() == float64(int64(t.Float())) {
			return int64(t.Float())
		}
		return 0
	case t.Kind() == reflect.String:
		i, _ := strconv.ParseFloat(t.String(), 64)
		return ToInt64(i)
	default:
		return 0
	}
}

func GetRole(g *discordgo.Guild, id string) *discordgo.Role {
	for i := range g.Roles {
		if g.Roles[i].ID == id {
			return g.Roles[i]
		}
	}

	return nil
}