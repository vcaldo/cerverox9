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

type VoiceMetrics struct {
	client influxdb2.Client
	org    string
	bucket string
	url    string
}

func NewAuthenticatedVoiceMetricsClient() *VoiceMetrics {
	return newVoiceMetricsClient(
		os.Getenv("INFLUX_URL"),
		os.Getenv("INFLUX_TOKEN"),
		os.Getenv("INFLUX_ORG"),
		os.Getenv("INFLUX_BUCKET"),
	)
}

func newVoiceMetricsClient(url, token, org, bucket string) *VoiceMetrics {
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	client := influxdb2.NewClient(url, token)
	return &VoiceMetrics{
		client: client,
		org:    org,
		bucket: bucket,
		url:    url,
	}
}

func (vm *VoiceMetrics) LogVoiceEvent(userID, username, channelID, eventType string, state bool) error {
	writeAPI := vm.client.WriteAPIBlocking(vm.org, vm.bucket)

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

func (vm *VoiceMetrics) LogOnlineUsers(guildID string, onlineUsers int) error {
	writeAPI := vm.client.WriteAPIBlocking(vm.org, vm.bucket)

	p := influxdb2.NewPoint(OnlineUsersMeasurement,
		map[string]string{
			"guild_id": guildID,
		},
		map[string]interface{}{
			"online_users": onlineUsers,
		},
		time.Now())
	log.Printf("Writing point: %s, %d in %s measurement", guildID, onlineUsers, OnlineUsersMeasurement)

	return writeAPI.WritePoint(context.Background(), p)
}

func (vm *VoiceMetrics) GetVoiceChatOnlineUsers(guildID string) (int, error) {
	query := fmt.Sprintf(`from(bucket:"%s")
		|> range(start: -10m)
		|> filter(fn: (r) => r._measurement == "%s" and r.guild_id == "%s")
		|> last()`,
		vm.bucket, OnlineUsersMeasurement, guildID)
	log.Println("Running query:", query)
	queryAPI := vm.client.QueryAPI(vm.org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return 0, fmt.Errorf("error querying for online users: %v", err)
	}
	defer result.Close()

	if result.Next() {
		record := result.Record()
		log.Printf("Record: %v", record.Values())
		onlineUsers := record.Values()[OnlineUsersMeasurement].(int)
		return onlineUsers, nil
	}
	return 0, fmt.Errorf("no online users found for guild %s", guildID)
}
