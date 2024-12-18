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
	case vsu.BeforeUpdate == nil && vsu.ChannelID != "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has joined voice channel %s", user.Username, vsu.ChannelID)
	case vsu.ChannelID == "":
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has left voice channel %s", user.Username, vsu.BeforeUpdate.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "leave", false)
	case vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != vsu.ChannelID:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has switched from voice channel %s to %s", user.Username, vsu.BeforeUpdate.ChannelID, vsu.ChannelID)
		// influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.BeforeUpdate.ChannelID, "leave", false)
	case !vsu.BeforeUpdate.SelfStream && vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has started streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", true)
	case vsu.BeforeUpdate.SelfStream && !vsu.SelfStream:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has stopped streaming in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "streaming", false)
	case !vsu.BeforeUpdate.SelfVideo && vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has turned on their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "video", true)
	case vsu.BeforeUpdate.SelfVideo && !vsu.SelfVideo:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has turned off their webcam in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "video", false)
	case !vsu.BeforeUpdate.SelfMute && vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has muted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", true)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "join", true)
		err = UpdateUsersInVoiceChannels(s)
		if err != nil {
			log.Println("error counting users in voice channels,", err)
			return
		}
	case vsu.BeforeUpdate.SelfMute && !vsu.SelfMute:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has unmuted themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "mute", false)
	case !vsu.BeforeUpdate.SelfDeaf && vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has deafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deaf", true)
	case vsu.BeforeUpdate.SelfDeaf && !vsu.SelfDeaf:
		user, err := s.User(vsu.UserID)
		if err != nil {
			log.Println("error fetching user:", err)
			return
		}
		log.Printf("User %s has undeafened themselves in voice channel %s", user.Username, vsu.ChannelID)
		influx.LogVoiceEvent(vsu.UserID, user.Username, vsu.ChannelID, "deaf", false)
	}
}

func UpdateUsersInVoiceChannels(s *discordgo.Session) error {
	influx := models.NewAuthenticatedVoiceMetricsClient()

	guilds, err := s.UserGuilds(200, "", "", true)
	if err != nil {
		return fmt.Errorf("error fetching guilds: %v", err)
	}
	for _, guild := range guilds {
		log.Printf("Guild: %s", guild.Name)
		guildID := guild.ID
		members, err := s.GuildMembers(guildID, "", 1000)
		if err != nil {
			log.Printf("error fetching members for guild %s: %v", guildID, err)
			continue
		}
		totalUsers := 0
		for _, member := range members {
			vs, _ := s.State.VoiceState(guildID, member.User.ID) // it errors out if the user is not in a voice channel, ignore it
			if vs != nil && vs.ChannelID != "" {
				totalUsers++
			}
		}

		err = influx.LogOnlineUsers(guildID, totalUsers)
		if err != nil {
			return fmt.Errorf("error logging online users: %v", err)
		}
	}
	return nil
}
