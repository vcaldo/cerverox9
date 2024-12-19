package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	VoiceEventsMeasurement = "voice_events"
	OnlineUsersMeasurement = "online_users"
	UserIdKey              = "user_id"
	UsernameKey            = "username"
	ChannelIdKey           = "channel_id"
	EventTypeKey           = "event_type"
	StateKey               = "state"
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

func (dm *DiscordMetrics) LogVoiceEvent(userID, username, channelID, eventType string, state bool) error {
	writeAPI := dm.Client.WriteAPIBlocking(dm.Org, dm.Bucket)

	p := influxdb2.NewPoint(VoiceEventsMeasurement,
		map[string]string{
			UserIdKey:    userID,
			UsernameKey:  username,
			ChannelIdKey: channelID,
			EventTypeKey: eventType,
		},
		map[string]interface{}{
			StateKey: state,
		},
		time.Now())
	log.Printf("Writing point: %s, %s, %t in %s measurement", username, eventType, state, VoiceEventsMeasurement)

	return writeAPI.WritePoint(context.Background(), p)
}

func (dm *DiscordMetrics) LogOnlineUsers(guildID string, onlineUsers int, userList []string) error {
	writeAPI := dm.client.WriteAPIBlocking(dm.Org, dm.Bucket)

	p := influxdb2.NewPoint(OnlineUsersMeasurement,
		map[string]string{
			"guild_id":  guildID,
			"user_list": strings.Join(userList, ","),
		},
		map[string]interface{}{
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
	log.Println("Running query:", query)
	queryAPI := dm.client.QueryAPI(dm.Org)
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
