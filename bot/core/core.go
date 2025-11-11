package core

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/RhykerWells/Summit/bot/events"
	"github.com/bwmarrin/discordgo"
)

// Init registers the required guild join & leave functions as well as initialises the web plugin
func Init() {
	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate){
		guildAddCoreConfig,
	})
	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete){
		guildDeleteCoreConfig,
	})

	initWeb()
}
