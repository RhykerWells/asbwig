package events

import "github.com/bwmarrin/discordgo"

// scheduledChannelCreateFunctions serves as a map of all the functions that is run when a guild creates a channel
var scheduledChannelCreateFunctions []func(g *discordgo.ChannelCreate)

// RegisterChannelCreatefunctions adds each guild channel create function to the map of functions ran when a guild creates a channel
func RegisterChannelCreatefunctions(funcMap []func(g *discordgo.ChannelCreate)) {
	scheduledChannelCreateFunctions = append(scheduledChannelCreateFunctions, funcMap...)
}

// channelCreate is called when a guild creates a channel
func channelCreate(s *discordgo.Session, c *discordgo.ChannelCreate) {
	for _, channelFunction := range scheduledChannelCreateFunctions {
		go channelFunction(c)
	}
}

// scheduledChannelUpdateFunctions serves as a map of all the functions that is run when a guild updates a channel
var scheduledChannelUpdateFunctions []func(g *discordgo.ChannelUpdate)

// RegisterChannelUpdatefunctions adds each guild channel update function to the map of functions ran when a guild updates a channel
func RegisterChannelUpdatefunctions(funcMap []func(g *discordgo.ChannelUpdate)) {
	scheduledChannelUpdateFunctions = append(scheduledChannelUpdateFunctions, funcMap...)
}

// channelUpdate is called when a guild updates a channel
func channelUpdate(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	for _, channelFunction := range scheduledChannelUpdateFunctions {
		go channelFunction(c)
	}
}

// scheduledChannelDeleteFunctions serves as a map of all the functions that is run when a guild deletes a channel
var scheduledChannelDeleteFunctions []func(g *discordgo.ChannelDelete)

// RegisterChannelDeletefunctions adds each guild channel delete function to the map of functions ran when a guild deletes a channel
func RegisterChannelDeletefunctions(funcMap []func(g *discordgo.ChannelDelete)) {
	scheduledChannelDeleteFunctions = append(scheduledChannelDeleteFunctions, funcMap...)
}

// channelDelete is called when a guild deletes a channel
func channelDelete(s *discordgo.Session, c *discordgo.ChannelDelete) {
	for _, channelFunction := range scheduledChannelDeleteFunctions {
		go channelFunction(c)
	}
}
