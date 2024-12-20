package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

func GetVoiceCallStatus() (oncallUsersCount int64, oncallUsers string, onlineUsersCount int64, onlineUsers string, error error) {
	dm := models.NewAuthenticatedDiscordMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")

	oncallUsersCount, oncallUsers, err := dm.GetOncallUsers(guildID)
	if err != nil {
		return 0, "", 0, "", err
	}

	onlineUsersCount, onlineUsers, err = dm.GetOnlineUsers(guildID)
	if err != nil {
		return 0, "", 0, "", err
	}

	return oncallUsersCount, oncallUsers, onlineUsersCount, onlineUsers, nil
}
