package handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

func VoiceEventHanlder(ctx context.Context, b *bot.Bot, event *VoiceEvent) {
	chatId, ok := os.LookupEnv("TELEGRAM_CHAT_ID")
	if !ok {
		panic("TELEGRAM_CHAT_IDenv var is required")
	}

	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		panic("TELEGRAM_CHAT_ID must be a valid int64")
	}

	switch {
	// User joined the voice channel
	case event.EventType == "voice" && event.State:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatIdInt,
			Text:   fmt.Sprintf("User %s entrou na festa online", event.Username),
		})
	// User left the voice channel
	case event.EventType == "voice" && !event.State:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatIdInt,
			Text:   fmt.Sprintf("User %s saiu da festa online", event.Username),
		})
	case event.EventType == "video" && event.State:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatIdInt,
			Text:   fmt.Sprintf("User %s abriu a cam na festa online", event.Username),
		})
	case event.EventType == "streaming" && event.State:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatIdInt,
			Text:   fmt.Sprintf("User %s começou a streamar na festa online ", event.Username),
		})
	}
}
