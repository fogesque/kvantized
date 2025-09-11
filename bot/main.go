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

	// Get schedule details from environment variables
	pWp := os.Getenv("P_WP")
	aWpEven := os.Getenv("A_WP_EVEN")
	aWpOdd := os.Getenv("A_WP_ODD")
	kWpEven := os.Getenv("B_WP_EVEN")
	kWpOdd := os.Getenv("B_WP_ODD")

	if pWp == "" || aWpEven == "" || aWpOdd == "" || kWpEven == "" || kWpOdd == "" {
		log.Fatal("Missing required schedule environment variables")
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Get current week number (ISO 8601)
	year, week := time.Now().ISOWeek()

	// Determine the current week's schedule
	var aWP, kWP string
	if week%2 == 0 {
		aWP = aWpEven
		kWP = kWpEven
	} else {
		aWP = aWpOdd
		kWP = kWpOdd
	}

	// Format the message
	message := fmt.Sprintf(
		"📅 *Неделя %d, год %d*\n\n"+
			"График на эту неделю:\n"+
			"P - %s\n"+
			"A - %s\n"+
			"K - %s\n",
		week, year,
		pWp, aWP, kWP,
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
