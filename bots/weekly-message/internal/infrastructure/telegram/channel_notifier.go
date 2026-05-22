package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChannelNotifier struct {
	channelID string
	bot       *tgbotapi.BotAPI
}

func NewChannelNotifier(botToken, channelID string) (ChannelNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return ChannelNotifier{}, err
	}

	return ChannelNotifier{
		channelID: channelID,
		bot:       bot,
	}, nil
}

func (n ChannelNotifier) Send(_ context.Context, message string) (int, error) {
	msg := tgbotapi.NewMessageToChannel(n.channelID, message)

	sentMsg, err := n.bot.Send(msg)
	if err != nil {
		return 0, err
	}

	return sentMsg.MessageID, nil
}
