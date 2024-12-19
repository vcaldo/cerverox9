package handlers

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

func NewVoiceEventListener() *VoiceEventListener {
	metrics := models.NewAuthenticatedDiscordMetricsClient()
	return &VoiceEventListener{
		metrics:    metrics,
		notifyChan: make(chan VoiceEvent, 100),
	}
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

		// Debug log
		log.Printf("Record values: %+v", values)

		// Safe value extraction
		userID, ok1 := values["user_id"].(string)
		username, ok2 := values["username"].(string)
		channelID, ok3 := values["channel_id"].(string)
		eventType, ok4 := values["event_type"].(string)
		guildID, ok5 := values["guild_id"].(string)

		// Get state value safely
		var state bool
		if stateVal, exists := values["state"]; exists && stateVal != nil {
			if stateBool, ok := stateVal.(bool); ok {
				state = stateBool
			} else {
				log.Printf("Invalid state type for record: %+v", values)
				continue
			}
		}

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
			Timestamp: record.Time(),
			GuildID:   guildID,
		}
		events = append(events, event)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	return events, nil
}
