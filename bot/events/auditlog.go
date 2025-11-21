package events

import "github.com/bwmarrin/discordgo"

var scheduledAuditLogCreateFunctions []func(g *discordgo.GuildAuditLogEntryCreate)

func RegisterAuditLogCreateFunctions(funcmap []func(g *discordgo.GuildAuditLogEntryCreate)) {
	scheduledAuditLogCreateFunctions = append(scheduledAuditLogCreateFunctions, funcmap...)
}

func auditLogCreate(s *discordgo.Session, g *discordgo.GuildAuditLogEntryCreate) {
	for _, auditLogCreateFunction := range scheduledAuditLogCreateFunctions {
		go auditLogCreateFunction(g)
	}
}
