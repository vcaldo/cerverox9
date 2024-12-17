package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("BOT_DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(voiceStateUpdate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}

func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	switch {
	case vsu.BeforeUpdate == nil && vsu.ChannelID != "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has joined voice channel %s", user.Username, vsu.ChannelID)
	case vsu.ChannelID == "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has left voice channel %s", user.Username, vsu.BeforeUpdate.ChannelID)
	case vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != vsu.ChannelID:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has switched from voice channel %s to %s", user.Username, vsu.BeforeUpdate.ChannelID, vsu.ChannelID)
	case !vsu.BeforeUpdate.SelfStream && vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has started streaming in voice channel %s", user.Username, vsu.ChannelID)
	case vsu.BeforeUpdate.SelfStream && !vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has stopped streaming in voice channel %s", user.Username, vsu.ChannelID)
	case !vsu.BeforeUpdate.SelfVideo && vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has turned on their webcam in voice channel %s", user.Username, vsu.ChannelID)
	case vsu.BeforeUpdate.SelfVideo && !vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has turned off their webcam in voice channel %s", user.Username, vsu.ChannelID)
	case !vsu.BeforeUpdate.SelfMute && vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has muted themselves in voice channel %s", user.Username, vsu.ChannelID)
	case vsu.BeforeUpdate.SelfMute && !vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has unmuted themselves in voice channel %s", user.Username, vsu.ChannelID)
	case !vsu.BeforeUpdate.SelfDeaf && vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has deafened themselves in voice channel %s", user.Username, vsu.ChannelID)
	case vsu.BeforeUpdate.SelfDeaf && !vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has undeafened themselves in voice channel %s", user.Username, vsu.ChannelID)
	}
}
