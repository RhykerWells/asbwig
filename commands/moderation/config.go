package moderation

import (
	"context"

	"github.com/RhykerWells/asbwig/commands/moderation/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type Config struct {
	// General
	GuildID	string
	ModerationEnabled bool
	ModerationTriggerDeletionEnabled bool
	ModerationTriggerDeletionSeconds int64
	ModerationResponseDeletionEnabled bool
	ModerationResponseDeletionSeconds int64
	ModerationLogChannel string

	// Warn
	WarnRequiredRoles []string

	// Mutes/Unmute
	MuteRequiredRoles []string
	MuteRole string
	MuteManageRole bool
	MuteUpdateRoles []string

	// Kick
	KickRequiredRoles []string

	// Bans/Unbans
	BanRequiredRoles []string

	LastCaseID int64
}

func (c *Config) ConfigToSQLModel() *models.ModerationConfig {
	return &models.ModerationConfig{
		// General
		GuildID: c.GuildID,
		ModerationEnabled: c.ModerationEnabled,
		ModerationTriggerDeletionEnabled: c.ModerationTriggerDeletionEnabled,
		ModerationTriggerDeletionSeconds: c.ModerationTriggerDeletionSeconds,
		ModerationResponseDeletionEnabled: c.ModerationResponseDeletionEnabled,
		ModerationResponseDeletionSeconds: c.ModerationResponseDeletionSeconds,
		ModerationLogChannel:  c.ModerationLogChannel,

		/* Warn */
		WarnRequiredRoles: c.WarnRequiredRoles,

		/* Mute/Unmute */
		MuteRequiredRoles: c.MuteRequiredRoles,
		MuteRole: c.MuteRole,
		MuteManageRole: c.MuteManageRole,
		MuteUpdateRoles: c.MuteUpdateRoles,

		/* Kick */
		KickRequiredRoles: c.KickRequiredRoles,

		/* Ban/Unban */
		BanRequiredRoles: c.BanRequiredRoles,

		LastCaseID: c.LastCaseID,
	}
}

func ConfigFromModel(m *models.ModerationConfig) *Config {
	return &Config{
		GuildID: m.GuildID,
		ModerationEnabled: m.ModerationEnabled,
		ModerationTriggerDeletionEnabled: m.ModerationTriggerDeletionEnabled,
		ModerationTriggerDeletionSeconds: m.ModerationTriggerDeletionSeconds,
		ModerationResponseDeletionEnabled: m.ModerationResponseDeletionEnabled,
		ModerationResponseDeletionSeconds: m.ModerationResponseDeletionSeconds,

		ModerationLogChannel: m.ModerationLogChannel,

		/* Warn */
		WarnRequiredRoles: m.WarnRequiredRoles,

		/* Mute */
		MuteRequiredRoles: m.MuteRequiredRoles,
		MuteRole: m.MuteRole,
		MuteManageRole: m.MuteManageRole,
		MuteUpdateRoles: m.MuteUpdateRoles,

		/* Kick */
		KickRequiredRoles: m.KickRequiredRoles,

		/* Ban */
		BanRequiredRoles: m.BanRequiredRoles,

		LastCaseID: m.LastCaseID,
	}
}

func GetConfig(guildID string) *Config {
	model, err := models.FindModerationConfigG(context.Background(), guildID)
	if err == nil {
		return ConfigFromModel(model)
	}

	return &Config{
		GuildID: guildID,
	}
}

func SaveConfig(config *Config) error {
	err := config.ConfigToSQLModel().UpsertG(context.Background(), true, []string{"guild_id"}, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func getGuildCases(guildID string) models.ModerationCaseSlice {
	models, err := models.ModerationCases(qm.Where("guild_id = ?", guildID)).All(context.Background(), common.PQ)
	if err != nil {
		return nil
	}

	return models
}