package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func Run() {
	discord, err := discordgo.New("Bot " + os.Getenv("ASBWIG_TOKEN"))
	if err != nil {
		log.Fatalln("ERROR LOGGING IN\n", err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatalln("ERROR OPENING CONNECTION\n", err)
	}
	defer discord.Close()
	closeHandler(discord)
}

func closeHandler(discord *discordgo.Session) {
	log.Infoln("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    // Cleanly close down the Discord session.
	log.Infoln("Exiting now....")
    discord.Close()
}