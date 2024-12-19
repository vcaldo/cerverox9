package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vcaldo/cerverox9/telegram/pkg/handlers"
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

	// Start the voice event listener
	listener := handlers.NewVoiceEventListener()
	go func() {
		listener.Start(ctx)
	}()

	chatId, ok := os.LookupEnv("TELEGRAM_CHAT_ID")
	if !ok {
		panic("TELEGRAM_CHAT_IDenv var is required")
	}

	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		panic("TELEGRAM_CHAT_ID must be a valid int64")
	}

	// Listen for events from the voice channel and send messages
	go func() {
		for event := range listener.NotifyChan {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatIdInt,
				Text:   fmt.Sprintf("User %s joined the voice channel", event.Username),
			})
		}
	}()

	// Wait for the context to be done
	select {}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	switch {
	case update.Message != nil && update.Message.Text == "/status":
		handlers.StatusHandler(ctx, b, update)
	}
}
