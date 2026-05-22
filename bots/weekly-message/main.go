package main

import (
	"context"
	"log"
	"net/http"

	"kvantized-bot/internal/application/weeklymessage"
	"kvantized-bot/internal/infrastructure/config"
	"kvantized-bot/internal/infrastructure/telegram"
	"kvantized-bot/internal/infrastructure/trello"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	board := trello.NewBoardClient(http.DefaultClient, trello.Credentials{
		APIKey:  settings.TrelloAPIKey,
		Token:   settings.TrelloToken,
		BoardID: settings.TrelloBoardID,
	})

	notifier, err := telegram.NewChannelNotifier(settings.TelegramBotToken, settings.TelegramChannelID)
	if err != nil {
		log.Fatal(err)
	}

	useCase := weeklymessage.NewSendTaskSummary(board, notifier, weeklymessage.TaskSummaryFormatter{})

	messageID, err := useCase.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent successfully! Message ID: %d", messageID)
}
