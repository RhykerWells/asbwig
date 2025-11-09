package events

import "github.com/bwmarrin/discordgo"

// scheduledGuildRoleCreateFunctions serves as a map of all functions that run when a guild role is updates
var scheduledGuildRoleCreateFunctions []func(r *discordgo.GuildRoleCreate)

// RegisterGuildRoleCreatefunctions adds each guild role create function to the map of functions ran when a guild creates a role
func RegisterGuildRoleCreatefunctions(funcMap []func(g *discordgo.GuildRoleCreate)) {
	scheduledGuildRoleCreateFunctions = append(scheduledGuildRoleCreateFunctions, funcMap...)
}

// guildRoleCreate is called when a role is deleted from a guild
func guildRoleCreate(s *discordgo.Session, r *discordgo.GuildRoleCreate) {
	for _, createFunction := range scheduledGuildRoleCreateFunctions {
		go createFunction(r)
	}
}

// scheduledGuildRoleUpdateFunctions serves as map of all functions that run when a guild role is updates
var scheduledGuildRoleUpdateFunctions []func(r *discordgo.GuildRoleUpdate)

// RegisterGuildRoleUpdatefunctions adds each guild role update function to the map of functions ran when a guild updates a role
func RegisterGuildRoleUpdatefunctions(funcMap []func(g *discordgo.GuildRoleUpdate)) {
	scheduledGuildRoleUpdateFunctions = append(scheduledGuildRoleUpdateFunctions, funcMap...)
}

// guildRoleUpdate is called when a role is updated in a guild
func guildRoleUpdate(s *discordgo.Session, r *discordgo.GuildRoleUpdate) {
	for _, updateFunction := range scheduledGuildRoleUpdateFunctions {
		go updateFunction(r)
	}
}

// scheduledGuildRoleDeleteFunctions serves as map of all functions that run when a guild role is deleted
var scheduledGuildRoleDeleteFunctions []func(r *discordgo.GuildRoleDelete)

// RegisterGuildRoleDeletefunctions adds each guild role delete function to the map of functions ran when a guild deletes a role
func RegisterGuildRoleDeletefunctions(funcMap []func(g *discordgo.GuildRoleDelete)) {
	scheduledGuildRoleDeleteFunctions = append(scheduledGuildRoleDeleteFunctions, funcMap...)
}

// guildRoleDelete is called when a role is deleted from a guild
func guildRoleDelete(s *discordgo.Session, r *discordgo.GuildRoleDelete) {
	for _, deleteFunction := range scheduledGuildRoleDeleteFunctions {
		go deleteFunction(r)
	}
}
