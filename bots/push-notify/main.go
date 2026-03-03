package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Get credentials from environment variables
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	channelID := os.Getenv("TELEGRAM_CHANNEL_ID")

	if botToken == "" || channelID == "" {
		log.Fatal("Missing required environment variables")
	}

	// Get commit details from environment variables
	commitMessages := os.Getenv("COMMIT_MESSAGES")
	changedFiles := os.Getenv("CHANGED_FILES")
	pusher := os.Getenv("PUSHER_NAME")
	repoName := os.Getenv("REPO_NAME")

	if commitMessages == "" {
		log.Fatal("No commit messages provided")
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Build commit list
	commits := strings.Split(commitMessages, "\n")
	var commitList strings.Builder
	for _, c := range commits {
		c = strings.TrimSpace(c)
		if c != "" {
			commitList.WriteString(fmt.Sprintf("• %s\n", c))
		}
	}

	// Build changed files list
	var fileList strings.Builder
	if changedFiles != "" {
		files := strings.Split(changedFiles, "\n")
		for _, f := range files {
			f = strings.TrimSpace(f)
			if f != "" {
				fileList.WriteString(fmt.Sprintf("• %s\n", f))
			}
		}
	}

	// Format the message
	message := fmt.Sprintf(
		"🚀 *Kvantized Update Notification*\n\n"+
			"*Repository:* %s\n"+
			"*Author:* %s\n\n"+
			"*Commits:*\n%s\n"+
			"*Changed Files:*\n%s",
		escapeMarkdown(repoName),
		escapeMarkdown(pusher),
		escapeMarkdown(commitList.String()),
		escapeMarkdown(fileList.String()),
	)

	// Create message config
	msg := tgbotapi.NewMessageToChannel(channelID, message)
	msg.ParseMode = "Markdown"

	// Send the message
	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Push notification sent successfully! Message ID: %d", sentMsg.MessageID)
}

// escapeMarkdown escapes special characters for Telegram Markdown V1.
func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"`", "\\`",
	)
	return replacer.Replace(s)
}
