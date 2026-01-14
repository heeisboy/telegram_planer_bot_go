package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tg-bots/internal/bot"
	"tg-bots/internal/config"
	"tg-bots/internal/scheduler"
	"tg-bots/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	botAPI, err := botapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("bot init error: %v", err)
	}

	store := storage.NewMemoryStore()
	h := bot.NewHandler(botAPI, store, cfg.DefaultTZ)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go scheduler.Run(ctx, 20*time.Second, store, store, func(userID int64, text string) error {
		msg := botapi.NewMessage(userID, text)
		_, err := botAPI.Send(msg)
		return err
	})

	updates := botAPI.GetUpdatesChan(botapi.NewUpdate(0))
	go func() {
		<-ctx.Done()
		botAPI.StopReceivingUpdates()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			h.HandleUpdate(ctx, update)
		}
	}
}
