package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	VoiceEventsMeasurement = "voice_events"
	OnlineUsersMeasurement = "online_users"
	UserIdKey              = "user_id"
	UsernameKey            = "username"
	UserDisplayName        = "user_display_name"
	GuildIdKey             = "guild_id"
	ChannelIdKey           = "channel_id"
	ChannelNameKey         = "channel_name"
	EventTypeKey           = "event_type"
	StateKey               = "state"
	VoiceEvent             = "voice"
	MuteEvent              = "mute"
	DeafenEvent            = "deafen"
	WebcamEvent            = "webcam"
	StreamEvent            = "streaming"
)

type DiscordMetrics struct {
	Client influxdb2.Client
	Org    string
	Bucket string
	Url    string
}

func NewAuthenticatedDiscordMetricsClient() *DiscordMetrics {
	return newDiscordMetricsClient(
		os.Getenv("INFLUX_URL"),
		os.Getenv("INFLUX_TOKEN"),
		os.Getenv("INFLUX_ORG"),
		os.Getenv("INFLUX_BUCKET"),
	)
}

func newDiscordMetricsClient(url, token, org, bucket string) *DiscordMetrics {
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	client := influxdb2.NewClient(url, token)
	return &DiscordMetrics{
		Client: client,
		Org:    org,
		Bucket: bucket,
		Url:    url,
	}
}

func (dm *DiscordMetrics) LogVoiceEvent(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate, channelID, voiceEvent string, state bool) error {
	NewAuthenticatedDiscordMetricsClient()

	user, err := s.User(vsu.UserID)
	if err != nil {
		return fmt.Errorf("error fetching user: %v", err)
	}

	channel, err := s.Channel(channelID)
	if err != nil {

		log.Println("error fetching channel:", err)
		return fmt.Errorf("error fetching channel: %v", err)
	}

	dm.logVoiceEvent(vsu.UserID, user.Username, user.Username, vsu.GuildID, channelID, channel.Name, voiceEvent, state)
	return nil
}

func (dm *DiscordMetrics) logVoiceEvent(userID, username, UserDisplayName, guildID, channelID, channelName, eventType string, state bool) error {
	writeAPI := dm.Client.WriteAPIBlocking(dm.Org, dm.Bucket)

	p := influxdb2.NewPoint(VoiceEventsMeasurement,
		map[string]string{
			UserIdKey:       userID,
			UsernameKey:     username,
			UserDisplayName: UserDisplayName,
			GuildIdKey:      guildID,
			ChannelIdKey:    channelID,
			ChannelNameKey:  channelName,
			EventTypeKey:    eventType,
		},
		map[string]interface{}{
			StateKey: state,
		},
		time.Now())
	log.Printf("Writing point: %s, %s, %t in %s measurement", username, eventType, state, VoiceEventsMeasurement)

	return writeAPI.WritePoint(context.Background(), p)
}

// func (dm *DiscordMetrics) LogVoiceChannelUsers(s *discordgo.Session) error {
// 	guilds, err := s.UserGuilds(200, "", "", true)
// 	if err != nil {
// 		return fmt.Errorf("error fetching guilds: %v", err)
// 	}
// 	for _, guild := range guilds {
// 		guildID := guild.ID
// 		members, err := s.GuildMembers(guildID, "", 1000)
// 		if err != nil {
// 			log.Printf("error fetching members for guild %s: %v", guildID, err)
// 			continue
// 		}
// 		totalOncallUsers := 0
// 		onlineUsers := []string{}
// 		for _, member := range members {
// 			vs, _ := s.State.VoiceState(guildID, member.User.ID) // it errors out if the user is not in a voice channel, ignore it
// 			if vs != nil && vs.ChannelID != "" {
// 				totalOncallUsers++
// 				user, err := s.User(member.User.ID)
// 				if err != nil {
// 					log.Printf("error fetching user %s: %v", member.User.ID, err)
// 					continue
// 				}

// 				onlineUsers = append(onlineUsers, fmt.Sprintf("%s - %s", user.Username, user.GlobalName))
// 			}
// 		}

// 		err = dm.LogOnlineUsers(guildID, totalOncallUsers, onlineUsers)
// 		if err != nil {
// 			return fmt.Errorf("error logging online users: %v", err)
// 		}
// 		log.Printf("Logged %d users in voice channels for guild %s", totalOncallUsers, guildID)
// 	}
// 	return nil
// }

func (dm *DiscordMetrics) LogOnlineUsers(guildID string, oncallUsers, onlineUsers int, oncallUserList, onlineUserList []string) error {
	writeAPI := dm.Client.WriteAPIBlocking(dm.Org, dm.Bucket)

	p := influxdb2.NewPoint(OnlineUsersMeasurement,
		map[string]string{
			"guild_id":    guildID,
			"oncall_list": strings.Join(oncallUserList, ","),
			"online_list": strings.Join(onlineUserList, ","),
		},
		map[string]interface{}{
			"oncall_users": oncallUsers,
			"online_users": onlineUsers,
		},
		time.Now())
	log.Printf("Writing point: %s, %d in %s measurement", guildID, onlineUsers, OnlineUsersMeasurement)

	return writeAPI.WritePoint(context.Background(), p)
}

func (dm *DiscordMetrics) GetVoiceChatOnlineUsers(guildID string) (int64, string, error) {
	query := fmt.Sprintf(`from(bucket:"%s")
		|> range(start: -10m)
		|> filter(fn: (r) => r._measurement == "%s" and r.guild_id == "%s")
		|> group(columns: ["guild_id"])
		|> sort(columns: ["_time"], desc: true)
		|> limit(n: 1)
		|> last()`,

		dm.Bucket, OnlineUsersMeasurement, guildID)
	queryAPI := dm.Client.QueryAPI(dm.Org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, "", fmt.Errorf("error querying for online users: %v", err)
	}
	defer result.Close()

	for result.Next() {
		record := result.Record()
		onlineUsers := record.Value().(int64)
		usersList := record.Values()["user_list"].(string)
		return onlineUsers, usersList, nil
	}

	return 0, "", fmt.Errorf("no online users found for guild %s", guildID)
}

func (dm *DiscordMetrics) RegisterVoiceChannelUsers(s *discordgo.Session) error {
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
		totalOncallUsers := 0
		totalOnlineUsers := 0
		oncallUsers := []string{}
		onlineUsers := []string{}
		for _, member := range members {
			vs, err := s.State.VoiceState(guildID, member.User.ID) // it errors out if the user is not in a voice channel, ignore it
			if err != nil {
				totalOnlineUsers++
				onlineUsers = append(onlineUsers, member.DisplayName())
				if vs != nil && vs.ChannelID != "" {
					totalOncallUsers++
					oncallUsers = append(oncallUsers, member.DisplayName())
				}
			}

			err = dm.LogOnlineUsers(guildID, totalOncallUsers, totalOnlineUsers, oncallUsers, onlineUsers)
			if err != nil {
				return fmt.Errorf("error logging online users: %v", err)
			}
			log.Printf("Logged %d users in voice channels and %d online usersfor guild %s", totalOncallUsers, totalOnlineUsers, guildID)
		}
		return nil
	}
	return nil
}
