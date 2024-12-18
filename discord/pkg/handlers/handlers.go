package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

func VoiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	influx := models.NewAuthenticatedVoiceMetricsClient()
	switch {
	// User joined a voice channel
	case vsu.BeforeUpdate == nil && vsu.ChannelID != "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has joined voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "voice", true)
		// Update the number of users in voice channels when a user joins
		if err := RegisterVoiceChannelUsers(s); err != nil {
			log.Println("error counting users in voice channels,", err)
			return
		}
		// Update user list when a user joins
		RegisterVoiceChannelUsers(s)
	// User left a voice channel
	case vsu.ChannelID == "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has left voice channel %s", user.Username, vsu.BeforeUpdate.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "voice", false)
		// Update the number of users in voice channels when a user leaves
		if err := RegisterVoiceChannelUsers(s); err != nil {
			log.Println("error counting users in voice channels,", err)
			return
		}
		// Update user list when a user joins
		RegisterVoiceChannelUsers(s)
	// User switched voice channels
	case vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != vsu.ChannelID:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has switched from voice channel %s to %s", user.Username, vsu.BeforeUpdate.ChannelID, vsu.ChannelID)
		// When user swtiches channels, they leave the previous one and join the new one
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "voice", false)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "voice", true)
	// User started streaming
	case !vsu.BeforeUpdate.SelfStream && vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has started streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", true)
	// User stopped streaming
	case vsu.BeforeUpdate.SelfStream && !vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has stopped streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", false)
	// User turned their webcam on
	case !vsu.BeforeUpdate.SelfVideo && vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has turned on their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "webcam", true)
	// User turned their webcam off
	case vsu.BeforeUpdate.SelfVideo && !vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has turned off their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "webcam", false)
	// User muted themselves
	case !vsu.BeforeUpdate.SelfMute && vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has muted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", true)
	// User unmuted themselves
	case vsu.BeforeUpdate.SelfMute && !vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has unmuted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", false)
	// User deafened themselves
	case !vsu.BeforeUpdate.SelfDeaf && vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has deafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deafen", true)
	// User undeafened themselves
	case vsu.BeforeUpdate.SelfDeaf && !vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has undeafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deafen", false)
	}
}

func RegisterVoiceChannelUsers(s *discordgo.Session) error {
	influx := models.NewAuthenticatedVoiceMetricsClient()

	guilds, err := s.UserGuilds(200, "", "", true)
	if err != nil {
		return fmt.Errorf("error fetching guilds: %v", err)
	}
	for _, guild := range guilds {
		guildID := guild.ID
		members, err := s.GuildMembers(guildID, "", 1000)
		if err != nil {
			log.Printf("error fetching members for guild %s: %v", guildID, err)
			continue
		}
		totalUsers := 0
		onlineUsers := []string{}
		for _, member := range members {
			vs, _ := s.State.VoiceState(guildID, member.User.ID) // it errors out if the user is not in a voice channel, ignore it
			if vs != nil && vs.ChannelID != "" {
				totalUsers++
				user, err := s.User(member.User.ID)
				if err != nil {
					log.Printf("error fetching user %s: %v", member.User.ID, err)
					continue
				}

				onlineUsers = append(onlineUsers, fmt.Sprintf("%s - %s", user.Username, user.GlobalName))
			}
		}

		err = influx.LogOnlineUsers(guildID, totalUsers, onlineUsers)
		if err != nil {
			return fmt.Errorf("error logging online users: %v", err)
		}
		log.Printf("Logged %d users in voice channels for guild %s", totalUsers, guildID)
	}
	return nil
}
