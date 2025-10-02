package events

import (
	"github.com/bwmarrin/discordgo"
)

// scheduledGuildMemberJoinFunctions serves as a map of all the functions that is run when a user joins a guild the bot is on
var scheduledGuildMemberJoinFunctions []func(g *discordgo.GuildMemberAdd)

// RegisterGuildMemberJoinfunctions adds each guild member join function to the map of functions ran when a user joins a guild the bot is on
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

// scheduledGuildMemberLeaveFunctions serves as a map of all the functions that is run when a user leaves a guild the bot is on
var scheduledGuildMemberLeaveFunctions []func(g *discordgo.GuildMemberRemove)

// RegisterGuildMemberLeavefunctions adds each guild member leave function to the map of functions ran when a user leaves a guild the bot is on
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
