package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	VoiceEventsMeasurement = "voice_events"
	UserIdKey              = "user_id"
	UsernameKey            = "username"
	ChannelIdKey           = "channel_id"
	EventTypeKey           = "event_type"
	StateKey               = "state"
	DefaultOrg             = "discord_org"
)

type VoiceMetrics struct {
	client influxdb2.Client
	org    string
	bucket string
	url    string
}

func NewVoiceMetrics(url, token, org, bucket string) *VoiceMetrics {
	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}
	client := influxdb2.NewClient(url, token)
	return &VoiceMetrics{
		client: client,
		org:    DefaultOrg,
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
	log.Printf("Writing point: %s, %s, %t", username, eventType, state)

	return writeAPI.WritePoint(context.Background(), p)
}
