package config

import (
	"errors"
	"os"
	"time"
)

// Config holds runtime configuration for the bot.
type Config struct {
	BotToken  string
	DefaultTZ *time.Location
}

func Load() (Config, error) {
	var cfg Config
	cfg.BotToken = os.Getenv("BOT_TOKEN")
	if cfg.BotToken == "" {
		return cfg, errors.New("BOT_TOKEN is required")
	}

	cfg.DefaultTZ = time.UTC
	if tzName := os.Getenv("DEFAULT_TZ"); tzName != "" {
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			return cfg, err
		}
		cfg.DefaultTZ = loc
	}

	return cfg, nil
}
