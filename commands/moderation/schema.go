package moderation

var GuildModerationSchema = []string{`
CREATE TABLE IF NOT EXISTS moderation_config (
	-- General
	guild_id TEXT PRIMARY KEY,
	moderation_enabled BOOL DEFAULT FALSE NOT NULL,
	moderation_trigger_deletion_enabled BOOL DEFAULT FALSE NOT NULL,
	moderation_trigger_deletion_seconds BIGINT DEFAULT 0 NOT NULL,
	moderation_response_deletion_enabled BOOL DEFAULT FALSE NOT NULL,
	moderation_response_deletion_seconds BIGINT DEFAULT 0 NOT NULL,

	moderation_log_channel TEXT DEFAULT '' NOT NULL,

	-- Warn
	warn_required_roles TEXT[] DEFAULT '{}' NOT NULL,

	-- Mute/Unmute
	mute_required_roles TEXT[] DEFAULT '{}' NOT NULL,
	mute_role TEXT DEFAULT '' NOT NULL,
	mute_manage_role BOOL DEFAULT FALSE NOT NULL,
	mute_update_roles TEXT[] DEFAULT '{}' NOT NULL,


	-- Kick
	kick_required_roles TEXT[] DEFAULT '{}' NOT NULL,

	-- Ban/Unban
	ban_required_roles TEXT[] DEFAULT '{}' NOT NULL,


	last_case_id BIGINT DEFAULT 0 NOT NULL
);
`,`
CREATE TABLE IF NOT EXISTS moderation_mutes (
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    roles TEXT[] NOT NULL,
	unmute_at TIMESTAMP NOT NULL,
    PRIMARY KEY (guild_id, user_id),
    CONSTRAINT fk_moderation_config_roles_guild FOREIGN KEY (guild_id)
        REFERENCES moderation_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE TABLE IF NOT EXISTS moderation_bans (
    guild_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
	unban_at TIMESTAMP NOT NULL,
    PRIMARY KEY (guild_id, user_id),
    CONSTRAINT fk_moderation_config_roles_guild FOREIGN KEY (guild_id)
        REFERENCES moderation_config (guild_id) ON DELETE CASCADE
);
`,`
CREATE TABLE IF NOT EXISTS moderation_cases (
	case_id BIGINT NOT NULL,
	guild_id TEXT NOT NULL,
	staff_id TEXT NOT NULL,
	staff_username TEXT NOT NULL,
	offender_id TEXT NOT NULL,
	offender_username TEXT NOT NULL,
	reason TEXT,
	action TEXT NOT NULL,
	log_link TEXT NOT NULL,
    PRIMARY KEY (guild_id, case_id),
	CONSTRAINT fk_guild_cases FOREIGN KEY (guild_id)
		REFERENCES moderation_config (guild_id) ON DELETE CASCADE
);
`,
}
