package main

import (
	"context"
	"os"
	"os/signal"

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

	// Wait for the context to be done
	select {}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	switch {
	case update.Message != nil && update.Message.Text == "/status":
		handlers.StatusHandler(ctx, b, update)
	}
}
