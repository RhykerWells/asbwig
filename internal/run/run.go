package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ranger-4297/asbwig/internal"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	Session *discordgo.Session
)

func Init() {
	err := internal.Init()
	if err != nil {
		log.WithError(err).Fatal("Failed to start core")
	}
}

func Run() {
	shutdown()
}

func shutdown() {
	sc := make(chan os.Signal, 2)
	signal.Notify(sc, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	log.Infoln("Exiting now....")
	os.Exit(0)
}