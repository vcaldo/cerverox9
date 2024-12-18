package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

func GetVoiceCallStatus() (string, error) {
	influx := models.NewAuthenticatedVoiceMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")
	onlineUsers, err := influx.GetVoiceChatOnlineUsers(guildID)
	if err != nil {
		return "", err
	}
	return onlineUsers, nil
}
