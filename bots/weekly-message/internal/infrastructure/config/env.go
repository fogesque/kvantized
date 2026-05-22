package config

import (
	"errors"
	"os"
	"strings"
)

type Settings struct {
	TelegramBotToken  string
	TelegramChannelID string
	TrelloAPIKey      string
	TrelloToken       string
	TrelloBoardID     string
}

func Load() (Settings, error) {
	settings := Settings{
		TelegramBotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChannelID: os.Getenv("TELEGRAM_CHANNEL_ID"),
		TrelloAPIKey:      os.Getenv("TRELLO_API_KEY"),
		TrelloToken:       os.Getenv("TRELLO_TOKEN"),
		TrelloBoardID:     os.Getenv("TRELLO_BOARD_ID"),
	}

	var missing []string
	if settings.TelegramBotToken == "" {
		missing = append(missing, "TELEGRAM_BOT_TOKEN")
	}
	if settings.TelegramChannelID == "" {
		missing = append(missing, "TELEGRAM_CHANNEL_ID")
	}
	if settings.TrelloAPIKey == "" {
		missing = append(missing, "TRELLO_API_KEY")
	}
	if settings.TrelloToken == "" {
		missing = append(missing, "TRELLO_TOKEN")
	}
	if settings.TrelloBoardID == "" {
		missing = append(missing, "TRELLO_BOARD_ID")
	}
	if len(missing) > 0 {
		return Settings{}, errors.New("missing required environment variables: " + strings.Join(missing, ", "))
	}

	return settings, nil
}
