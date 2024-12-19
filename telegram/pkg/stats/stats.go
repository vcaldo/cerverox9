package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/telegram/pkg/models"
)

func GetVoiceCallStatus() (int64, string, error) {
	dm := models.NewAuthenticatedDiscordMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")
	onlineUsers, userList, err := dm.GetVoiceChatOnlineUsers(guildID)
	if err != nil {
		return 0, "", err
	}
	return onlineUsers, userList, nil
}
