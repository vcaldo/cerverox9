package stats

import (
	"os"

	"github.com/vcaldo/cerverox9/discord/pkg/models"
)

func GetVoiceCallStatus() (int64, string, error) {
	dm := models.NewAuthenticatedDiscordMetricsClient()
	guildID := os.Getenv("DISCORD_GUILD_ID")
	onlineUsers, userList, err := dm.GetGuildStats(guildID)
	if err != nil {
		return 0, "", err
	}
	return onlineUsers, userList, nil
}
