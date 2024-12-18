package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vcaldo/cerverox9/discord/pkg/handlers"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("DISCORD_BOT_TOKEN")

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	// Register necessary Intents for the bot
	dg.Identify.Intents = discordgo.IntentGuilds |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuildVoiceStates

	dg.AddHandler(handlers.VoiceStateUpdate)

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	log.Println("Discord Bot is now running.")

	// Launch a goroutine to update the number of users in voice channels when the bot starts
	go handlers.UpdateUsersInVoiceChannels(dg)

	// Update the number of users in voice channels every 60 seconds
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				handlers.UpdateUsersInVoiceChannels(dg)
			}
		}
	}()

	select {}
}
