package weeklymessage

import (
	"context"
	"fmt"

	"kvantized-bot/internal/domain/taskboard"
)

type BoardSummaryProvider interface {
	ActiveTaskSummary(ctx context.Context) (taskboard.TaskSummary, error)
}

type Notifier interface {
	Send(ctx context.Context, message string) (int, error)
}

type Formatter interface {
	Format(summary taskboard.TaskSummary) string
}

type SendTaskSummary struct {
	board     BoardSummaryProvider
	notifier  Notifier
	formatter Formatter
}

func NewSendTaskSummary(board BoardSummaryProvider, notifier Notifier, formatter Formatter) SendTaskSummary {
	return SendTaskSummary{
		board:     board,
		notifier:  notifier,
		formatter: formatter,
	}
}

func (uc SendTaskSummary) Run(ctx context.Context) (int, error) {
	summary, err := uc.board.ActiveTaskSummary(ctx)
	if err != nil {
		return 0, fmt.Errorf("get active task summary: %w", err)
	}

	messageID, err := uc.notifier.Send(ctx, uc.formatter.Format(summary))
	if err != nil {
		return 0, fmt.Errorf("send weekly message: %w", err)
	}

	return messageID, nil
}
