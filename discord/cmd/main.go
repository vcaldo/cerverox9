package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/vcaldo/cerverox9/discord/pkg/handlers"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(handlers.VoiceStateUpdate)

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	log.Println("Bot is now running.")
	select {}
}
