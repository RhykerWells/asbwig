package run

import (
	"fmt"
	"os"
	"os/signal"
	"log"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	discord, err := discordgo.New("Bot " + os.Getenv("ASBWIG_TOKEN"))
	if err != nil {
		log.Fatal("ERROR LOGGING IN\n", err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatal("ERROR OPENING CONNECTION\n", err)
	}
	defer discord.Close()
	closeHandler(discord)
}

func closeHandler(discord *discordgo.Session) {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    // Cleanly close down the Discord session.
	fmt.Println("Exiting now....")
    discord.Close()
}