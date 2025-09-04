package moderation

import (
	"context"
	"fmt"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/durationutil"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type logAction struct {
	CaseType string
	Name   string
	Colour int
}

var (
	logWarn   = logAction{Name: "Warned", CaseType: "Warning", Colour: 0xFCA253}
	logMute   = logAction{Name: "Muted", CaseType: "Mute",Colour: 0x5772BE}
	logUnmute = logAction{Name: "Unmuted", CaseType: "Unmute",Colour: common.SuccessGreen}
	logKick   = logAction{Name: "Kicked", CaseType: "Kick",Colour: 0xF2A013}
	logBan    = logAction{Name: "Banned", CaseType: "Ban",Colour: 0xD64848}
	logUnban  = logAction{Name: "Unbanned", CaseType: "Unban",Colour: common.SuccessGreen}
)

const (
	caseEmoji    string = "<:ID:1369739780958457966>"
	userEmoji    string = "<:Member:1369740929568739499>"
	actionEmoji  string = "<:Action:1369745870001799321>"
	channelEmoji string = "<:Channel:1369743815887294687>"
	reasonEmoji  string = "<:Reason:1369744280310124624>"
)

// Case generation

// setupModerationCase returns the models.ModerationCase struct of data to insert into the database
func caseUpsert(caseID int64, guildID, staffID, targetID, reason, loglink string, action logAction) error {
	moderationCase := models.ModerationCase{
		CaseID: caseID,
		GuildID: guildID,
		StaffID: staffID,
		OffenderID: targetID,
		Reason: null.StringFrom(reason),
		Action: action.Name,
		Loglink: loglink,
	}
	
	return addCase(guildID, moderationCase)
}

// getNewCaseID returns the caseID for the currentCase
func getNewCaseID(guildID string) int64 {
	guild, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	currentCaseID := guild.LastCaseID + 1
	return currentCaseID
}

// incrementCaseID increments the guilds caseID. This is only used once everything else is successful. Such as posting the log and setting up the case data.
func incrementCaseID(guildID string) error {
	guild, _ := models.ModerationConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	guild.LastCaseID++
	_, err := guild.Update(context.Background(), common.PQ, boil.Infer())
	return err
}

// addCase adds the case to the database and runs incrementCaseID.
func addCase(guildID string, caseData models.ModerationCase) error {
	err := caseData.Insert(context.Background(), common.PQ, boil.Infer())
	if err != nil {
		return err
	}
	err = incrementCaseID(guildID)
	if err != nil {
		removeFailedCase(caseData)
		return err
	}
	return nil
}

func removeFailedCase(caseData models.ModerationCase) {
	caseData.Delete(context.Background(), common.PQ)
}

// Log generation

// logCase runs the generation of the modlog embed and case upsertion. Returning an error if it wasn't complete
func logCase(guildID string, Author, target *discordgo.Member, action logAction, currentChannel, reason string, duration ...time.Duration) error {
	logChannel, err := getGuildModLogChannel(guildID)
	if err != nil {
		return err
	}
	caseID := getNewCaseID(guildID)
	embed := logEmbed(Author.User, target.User, caseID, action, currentChannel, reason, duration...)
	message, err := functions.SendMessage(logChannel, &discordgo.MessageSend{Embed: embed})
	if err != nil {
		return err
	}
	loglink := generateLogLink(guildID, logChannel, message.ID)
	caseUpsert(caseID, guildID, Author.User.ID, target.User.ID, reason, loglink, action)
	return nil
}

// logEmbed returns the fully-populated embed for moderation logging
func logEmbed(author, target *discordgo.User, caseNumber int64, action logAction, channelID, reason string, duration ...time.Duration) *discordgo.MessageEmbed {
	humanReadableCaseNumber := humanize.Comma(caseNumber)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s (ID %s)", author.Username, author.ID),
			IconURL: author.AvatarURL("1024"),
		},
		Description: fmt.Sprintf("%s **Case number:** %s\n%s **Who:** %s `(ID %s)`\n%s **Action:** %s\n%s **Channel:** <#%s>\n%s **Reason:** %s", caseEmoji, humanReadableCaseNumber, userEmoji, target.Mention(), target.ID, actionEmoji, action.Name, channelEmoji, channelID, reasonEmoji, reason),
		Color: action.Colour,
	}
	if len(duration) > 0 {
		d := duration[0]
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Duration: %s", durationutil.HumanizeDuration(d) ),
		}
	}
	return embed
}

// generateLogLink returns the full messageURL of the modlog entry
func generateLogLink(guildID, channelID, messageID string) string {
	link := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID)
	return link
}