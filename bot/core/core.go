package core

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/RhykerWells/asbwig/bot/events"
	"github.com/bwmarrin/discordgo"
)

func Init() {
	events.RegisterGuildJoinfunctions([]func(g *discordgo.GuildCreate) {
		guildAddCoreConfig,
	})
	events.RegisterGuildLeavefunctions([]func(g *discordgo.GuildDelete) {
		guildDeleteCoreConfig,
	})

	initWeb()
}