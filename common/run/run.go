package run

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RhykerWells/asbwig/bot"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/web"
	log "github.com/sirupsen/logrus"
)

func Init() {
	err := common.Init()
	bot.Run(common.Session)
	common.Run(common.Session)
	web.Run()
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
	shutdownCleanup()
	os.Exit(0)
}

func shutdownCleanup() {
	log.Warnln("Running cleanup functions")
	ok, err := models.EconomyCreateitems().DeleteAllG(context.Background())
	if err != nil {
		log.Errorln("Error running cleanup for economy:", ok)
	}
}