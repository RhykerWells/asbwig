package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/RhykerWells/asbwig/bot"
	"github.com/RhykerWells/asbwig/common"
	log "github.com/sirupsen/logrus"
)

func Init() {
	err := common.Init()
	bot.Run(common.Session)
	common.Run(common.Session)
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
