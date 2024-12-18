package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vcaldo/cerverox9/telegram/pkg/stats"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	// Start the bot in a goroutine
	go func() {
		b.Start(ctx)
	}()

	// Wait for the context to be done
	select {}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	switch {
	case update.Message != nil && update.Message.Text == "/status":
		onlineUsers, err := stats.GetVoiceCallStatus()
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erro ao buscar status da festa online",
			})
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%d users estão se divertindo na festa online", onlineUsers),
		})
		return
	}
}
