package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/telegram/pkg/models"
)

func GetVoiceCallStatus() (int, error) {
	influx := models.NewAuthenticatedVoiceMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")
	onlineUsers, err := influx.GetVoiceChatOnlineUsers(guildID)
	if err != nil {
		return 0, err
	}
	return onlineUsers, nil
}
