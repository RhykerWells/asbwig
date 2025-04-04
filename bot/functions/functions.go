package functions

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/RhykerWells/asbwig/common"
	"github.com/bwmarrin/discordgo"
)

// Message functions

// SendBasicMessage sends a string as message content to the given channel
// If a delay is included, then the message is deleted after X seconds.
func SendBasicMessage(channelID string, message string, delay ...any) (msg *discordgo.Message, err error) {
	msg, err = common.Session.ChannelMessageSend(channelID, message)
	if delay != nil {
		DeleteMessage(channelID, msg.ID, delay)
	}
	return msg, err
}

// SendMessage sends complex message objects to a given channel. Supporting, embed, components etc.
// If a delay is included, then the message is deleted after X seconds
func SendMessage(channelID string, message *discordgo.MessageSend, delay ...any) (msg *discordgo.Message, err error) {
	msg, err = common.Session.ChannelMessageSendComplex(channelID, message)
	if delay != nil {
		DeleteMessage(channelID, msg.ID, delay)
	}
	return msg, err
}

// SendDM sends complex message objects to a given users DM channel. Supporting, embed, components etc.
func SendDM(userID string, message *discordgo.MessageSend) error {
	channel, err := common.Session.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = SendMessage(channel.ID, message)

	return err
}

// EditBasicMessage edits a 'basic' message and allows replacement of the message content
func EditBasicMessage(channelID, messageID, message string) error {
	_, err := common.Session.ChannelMessageEdit(channelID, messageID, message)
	return err
}

// EditMessage edits a 'complex' message and allows replacement of all message objects
func EditMessage(channelID string, messageID string, message *discordgo.MessageSend) error {
	edit := &discordgo.MessageEdit{
		ID:	messageID,
		Channel: channelID,
	}
	if message.Content != "" {edit.Content = &message.Content}
	if message.Embed != nil {edit.Embed = message.Embed}
	if message.Embeds != nil {edit.Embeds = &message.Embeds}
	if message.Components != nil {edit.Components = &message.Components}
	_, err := common.Session.ChannelMessageEditComplex(edit)

	return err
}

// DeleteMessage deletes a given message after 0 or an option delay
func DeleteMessage(channelID, messageID string, delay ...any) error {
	var duration int
	if len(delay) > 0 {
		duration = delay[0].([]any)[0].(int)
    }
	time.Sleep(time.Duration(duration * int(time.Second)))
	err := common.Session.ChannelMessageDelete(channelID, messageID)

	return err
}

// User functions

// GetUser returns the user object if possible of a user ID
func GetUser(userID string) (*discordgo.User, error) {
	u, err := common.Session.User(userID)

	return u, err
}

// GetUser returns the member object if possible of a user ID
func GetMember(guildID string, userID string) (*discordgo.Member, error) {
	// Direct mention
	if strings.HasPrefix(userID, "<@") {
		userID = userID[2 : len(userID)-1]
	}
	u, err := common.Session.GuildMember(userID, userID)

	return u, err
}

// Role functions

// GetRole returns the full guild role object for a role ID/mention
func GetRole(guildID, roleStr string) (role *discordgo.Role, err error) {
	guild, err := common.Session.Guild(guildID)
	if err != nil {
		return nil, err
	}
	// Role mention
	if strings.HasPrefix(roleStr, "<@") {
		roleStr = roleStr[3 : len(roleStr)-1]
	}
	for i := range guild.Roles {
		if guild.Roles[i].ID == roleStr {
			role = guild.Roles[i]
			break
		}
	}
	return role, nil
}

// AddRole adds a given roleID to a user
func AddRole(guildID string, memberID, roleID string) error {
	member, err := GetMember(guildID, memberID)
	if err != nil {
		return err
	}
	if slices.Contains(member.Roles, roleID) {
		// User has the role
		return nil
	}
	err = common.Session.GuildMemberRoleAdd(guildID, memberID, roleID)

	return err
}

// AddRole removes a given roleID to a user
func RemoveRole(guildID, memberID, roleID string) error {
	member, err := GetMember(guildID, memberID)
	if err != nil {
		return err
	}
	if !slices.Contains(member.Roles, roleID) {
		// User doesn't have role
		return err
	}
	err = common.Session.GuildMemberRoleRemove(guildID, memberID, roleID)

	return err
}

// SetRoles flushes a users roles and reassigns them to the given roleIDs
func SetRoles(guildID, memberID string, roleIDs []string) error {
	member, err := GetMember(guildID, memberID)
	if err != nil {
		return err
	}
	roles := make(map[string]struct{})

	for _, id := range member.Roles {
		role, _ := GetRole(guildID, id)
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
	_, err = common.Session.GuildMemberEdit(guildID, memberID, userData)

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

// ToInt64 takes the value of an int, float or string and returns it as a whole 64-bit integer if possible or 0 when not.
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