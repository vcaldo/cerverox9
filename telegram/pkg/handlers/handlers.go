package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vcaldo/cerverox9/telegram/pkg/stats"
)

func StatusHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	onlineUsers, usersList, err := stats.GetVoiceCallStatus()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Erro ao buscar status da festa online",
		})
		return
	}
	userSlice := strings.Split(usersList, ",")
	usersListLineBreak := strings.Join(userSlice, "\n")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("%d users estão se divertindo na festa online 🥳\n\nUsers na festa online:\n%s", onlineUsers, usersListLineBreak),
	})
}
