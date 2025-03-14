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
func SendMessage(c string, messageData *discordgo.MessageSend, delay ...any) error {
	message, err := internal.Session.ChannelMessageSendComplex(c, messageData)
	DeleteMessage(c, message.ID, delay)
	return err
}

func SendDM(user string, msg *discordgo.MessageSend) error {
	channel, err := internal.Session.UserChannelCreate(user)
	if err != nil {
		return err
	}

	err = SendMessage(channel.ID, msg)
	return err
}

func DeleteMessage(c, m string, delay ...any) error {
    dur := 0
    if len(delay) > 0 {
        dur = int(ToInt64(delay[0]))
    }
	time.Sleep(time.Duration(dur))

    logrus.Infoln(dur)
    err := internal.Session.ChannelMessageDelete(c, m)
    return err
}

// User functions
func GetUser(u string) (interface{}, error) {
	user, err := internal.Session.User(u)

	return user, err
}

func GetMember(g *discordgo.Guild, u string) (interface{}, error) {
	user, err := internal.Session.GuildMember(g.ID, u)

	return user, err
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