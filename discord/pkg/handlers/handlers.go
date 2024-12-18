package handlers

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

func VoiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	influx := models.NewVoiceMetrics(
		os.Getenv("INFLUX_URL"),
		os.Getenv("INFLUX_TOKEN"),
		os.Getenv("INFLUX_ORG"),
		os.Getenv("INFLUX_BUCKET"),
	)
	switch {
	case vsu.BeforeUpdate == nil && vsu.ChannelID != "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has joined voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "join", true)
	case vsu.ChannelID == "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has left voice channel %s", user.Username, vsu.BeforeUpdate.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "leave", false)
	case vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != vsu.ChannelID:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has switched from voice channel %s to %s", user.Username, vsu.BeforeUpdate.ChannelID, vsu.ChannelID)
		// influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "leave", false)
	case !vsu.BeforeUpdate.SelfStream && vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has started streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", true)
	case vsu.BeforeUpdate.SelfStream && !vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has stopped streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", false)
	case !vsu.BeforeUpdate.SelfVideo && vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has turned on their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "video", true)
	case vsu.BeforeUpdate.SelfVideo && !vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has turned off their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "video", false)
	case !vsu.BeforeUpdate.SelfMute && vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has muted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", true)
	case vsu.BeforeUpdate.SelfMute && !vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has unmuted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", false)
	case !vsu.BeforeUpdate.SelfDeaf && vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has deafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deaf", true)
	case vsu.BeforeUpdate.SelfDeaf && !vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("Error fetching user:", err)
			return
		}
		log.Printf("User %s has undeafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deaf", false)
	}
}
