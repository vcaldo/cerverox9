package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/telegram/pkg/models"
)

func GetVoiceCallStatus() (int64, string, error) {
	influx := models.NewAuthenticatedVoiceMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")
	onlineUsers, userList, err := influx.GetVoiceChatOnlineUsers(guildID)
	if err != nil {
		return 0, "", err
	}
	return onlineUsers, userList, nil
}
