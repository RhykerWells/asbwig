package bot

import (
	"reflect"
	"strconv"
	"time"

	"github.com/Ranger-4297/asbwig/internal"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Message functions
func SendMessage(channelID string, messageData *discordgo.MessageSend, delay ...any) error {
	message, err := internal.Session.ChannelMessageSendComplex(channelID, messageData)
	DeleteMessage(channelID, message.ID, delay)
	return err
}

func SendDM(userID string, messageData *discordgo.MessageSend) error {
	channel, err := internal.Session.UserChannelCreate(userID)
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
    err := internal.Session.ChannelMessageDelete(channelID, messageData)
    return err
}

// User functions
func GetUser(user string) (interface{}, error) {
	u, err := internal.Session.User(user)

	return u, err
}

func GetMember(guild *discordgo.Guild, user string) (interface{}, error) {
	u, err := internal.Session.GuildMember(guild.ID, user)

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

    return internal.Session.GuildMemberRoleAdd(guild.ID, member.User.ID, roleID)
}

func RemoveRole(guild *discordgo.Guild, member *discordgo.Member, roleID string) error {
	for _, v := range member.Roles {

        if GetRole(guild, v).ID != roleID {
            internal.Session.GuildMemberRoleRemove(guild.ID, member.User.ID, roleID)
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
	_, err := internal.Session.GuildMemberEdit(guild.ID, member.User.ID, userData)
	return err
}

// Helper tools
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