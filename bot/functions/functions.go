package functions

import (
	"reflect"
	"strconv"
	"time"

	"github.com/Ranger-4297/asbwig/common"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Message functions
func SendBasicMessage(channelID string, message string) error {
	_, err := common.Session.ChannelMessageSend(channelID, message)
	return err
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

func DeleteMessage(channelID, messageData string, delay ...any) error {
	duration := 0
	if len(delay) > 0 {
		if int(ToInt64(delay[0])) < 1 {
			return nil
		}
		duration = int(ToInt64(delay[0]))
	}
	time.Sleep(time.Duration(duration))

	logrus.Infoln(duration)
	err := common.Session.ChannelMessageDelete(channelID, messageData)
	return err
}

// User functions
func GetUser(user string) (interface{}, error) {
	u, err := common.Session.User(user)

	return u, err
}

func GetMember(guild *discordgo.Guild, user string) (interface{}, error) {
	u, err := common.Session.GuildMember(guild.ID, user)

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
func ToInt64(conv interface{}) int64 {
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
