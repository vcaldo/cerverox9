package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

type VoiceEvent struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	ChannelID string `json:"channel_id"`
	EventType string `json:"event_type"`
	State     bool   `json:"state"`
}

type VoiceEventListener struct {
	Metrics     *models.DiscordMetrics
	LastChecked time.Time
	NotifyChan  chan VoiceEvent
}

func NewVoiceEventListener() *VoiceEventListener {
	metrics := models.NewAuthenticatedDiscordMetricsClient()
	return &VoiceEventListener{
		Metrics:    metrics,
		NotifyChan: make(chan VoiceEvent, 200),
	}
}

func (l *VoiceEventListener) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(l.NotifyChan)
			return
		case <-ticker.C:
			events, err := l.checkNewEvents()
			if err != nil {
				log.Printf("Error checking events: %v", err)
				continue
			}
			for _, event := range events {
				select {
				case l.NotifyChan <- event:
				default:
					log.Println("Channel buffer full, skipping event")
				}
			}
			l.LastChecked = time.Now()
		}
	}
}

func (l *VoiceEventListener) NotificationChannel() <-chan VoiceEvent {
	return l.NotifyChan
}

func (l *VoiceEventListener) checkNewEvents() ([]VoiceEvent, error) {
	query := fmt.Sprintf(`from(bucket:"%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r._measurement == "voice_events" and (r.event_type == "voice" or r.event_type == "webcam" or r.event_type == "streaming"))
		|> sort(columns: ["_time"])`,
		l.Metrics.Bucket,
		l.LastChecked.Format(time.RFC3339),
		time.Now().Format(time.RFC3339))

	result, err := l.Metrics.Client.QueryAPI(l.Metrics.Org).Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer result.Close()

	var events []VoiceEvent
	for result.Next() {
		record := result.Record()
		values := record.Values()

		// Safe value extraction
		userID, ok1 := values["user_id"].(string)
		username, ok2 := values["username"].(string)
		channelID, ok3 := values["channel_id"].(string)
		eventType, ok4 := values["event_type"].(string)
		state, ok5 := record.Value().(bool)

		// Skip if required fields are missing
		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
			log.Printf("Skipping record with missing fields: %+v", values)
			continue
		}

		event := VoiceEvent{
			UserID:    userID,
			Username:  username,
			ChannelID: channelID,
			EventType: eventType,
			State:     state,
		}
		events = append(events, event)
		log.Printf("Event: %+v", event)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return events, nil
}
