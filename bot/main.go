package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Get credentials from environment variables
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	channelID := os.Getenv("TELEGRAM_CHANNEL_ID")

	if botToken == "" || channelID == "" {
		log.Fatal("Missing required environment variables")
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Get current week number (ISO 8601)
	year, week := time.Now().ISOWeek()

	// Format the message
	message := fmt.Sprintf(
		"📅 *Неделя %d, год %d*\n\n"+
			"_это_автоматическое_сообщение_",
		week, year,
	)

	// Create message config
	msg := tgbotapi.NewMessageToChannel(channelID, message)
	msg.ParseMode = "Markdown"

	// Send the message
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent successfully! Message ID: %d", sentMsg.MessageID)
}
