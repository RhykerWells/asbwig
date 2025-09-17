package events

import (
	"github.com/bwmarrin/discordgo"
)

var scheduledGuildMemberJoinFunctions []func(g *discordgo.GuildMemberAdd)

func RegisterGuildMemberJoinfunctions(funcMap []func(g *discordgo.GuildMemberAdd)) {
	scheduledGuildMemberJoinFunctions = append(scheduledGuildMemberJoinFunctions, funcMap...)
}


// guildMemberAdd is called when a member joins a guild the bot is in
// This adds the user to any tables that are relevant to them
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	for _, joinFunction := range scheduledGuildMemberJoinFunctions {
		joinFunction(m)
	}
}

var scheduledGuildMemberLeaveFunctions []func(g *discordgo.GuildMemberRemove)

func RegisterGuildMemberLeavefunctions(funcMap []func(g *discordgo.GuildMemberRemove)) {
	scheduledGuildMemberLeaveFunctions = append(scheduledGuildMemberLeaveFunctions, funcMap...)
}

// guildMemberLeave is called when a member leaves a guild the bot is in
// This removes the user from any tables that they may be part of
func guildMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	for _, leaveFunction := range scheduledGuildMemberLeaveFunctions {
		leaveFunction(m)
	}
}
