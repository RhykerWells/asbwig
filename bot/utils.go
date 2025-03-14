package bot

import (
	"time"

	"github.com/Ranger-4297/asbwig/internal"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Message functions
func SendMessage(c string, messageData *discordgo.MessageSend) (messageObject *discordgo.Message, err error) {
	messageObject, err = internal.Session.ChannelMessageSendComplex(c, messageData)
	return
}

func SendDM(user string, msg *discordgo.MessageSend) error {
	channel, err := internal.Session.UserChannelCreate(user)
	if err != nil {
		return err
	}

	_, err = SendMessage(channel.ID, msg)
	return err
}

func SendMessageToDeleteAfter(delay time.Duration, c string, messageData *discordgo.MessageSend) {
	message, err := SendMessage(c, messageData)
	if err != nil {
		logrus.Warn("Failed to send temporary message")
	}
	DeleteChannelMessageAfterDelay(delay, c, message.ID)
}

func DeleteChannelMessageAfterDelay(delay time.Duration, c string, m string) error {
    time.Sleep(delay)
    err := internal.Session.ChannelMessageDelete(c, m)
    return err
}

// User functions
func GetUser(u string) (interface{}, error) {
	user, err := internal.Session.User(u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetMember(g, u string) (interface{}, error) {
	user, err := internal.Session.GuildMember(g, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}