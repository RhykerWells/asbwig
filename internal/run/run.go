package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ranger-4297/asbwig/bot"
	"github.com/Ranger-4297/asbwig/internal"
	log "github.com/sirupsen/logrus"
)

func Init() {
	err := internal.Init()
	bot.Run(internal.Session)
	internal.Run(internal.Session)
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