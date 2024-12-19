package voiceevents

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

type VoiceEvent struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	ChannelID string    `json:"channel_id"`
	EventType string    `json:"event_type"`
	State     bool      `json:"state"`
	Timestamp time.Time `json:"timestamp"`
	GuildID   string    `json:"guild_id"`
}

type VoiceEventListener struct {
	metrics     *models.DiscordMetrics
	lastChecked time.Time
	notifyChan  chan VoiceEvent
}

func (l *VoiceEventListener) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(l.notifyChan)
			return
		case <-ticker.C:
			events, err := l.checkNewEvents()
			if err != nil {
				log.Printf("Error checking events: %v", err)
				continue
			}
			for _, event := range events {
				select {
				case l.notifyChan <- event:
				default:
					log.Println("Channel buffer full, skipping event")
				}
			}
			l.lastChecked = time.Now()
		}
	}
}

func (l *VoiceEventListener) NotificationChannel() <-chan VoiceEvent {
	return l.notifyChan
}

func (l *VoiceEventListener) checkNewEvents() ([]VoiceEvent, error) {
	query := fmt.Sprintf(`
        from(bucket:"%s")
            |> range(start: %s)
            |> filter(fn: (r) => r._measurement == "voice_events")
            |> sort(columns: ["_time"])`,
		l.metrics.Bucket,
		l.lastChecked.Format(time.RFC3339))

	result, err := l.metrics.Client.QueryAPI(l.metrics.Org).Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer result.Close()

	var events []VoiceEvent
	for result.Next() {
		record := result.Record()
		values := record.Values()

		event := VoiceEvent{
			UserID:    values["user_id"].(string),
			Username:  values["username"].(string),
			ChannelID: values["channel_id"].(string),
			EventType: values["event_type"].(string),
			State:     values["state"].(bool),
			Timestamp: record.Time(),
		}
		events = append(events, event)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return events, nil
}
