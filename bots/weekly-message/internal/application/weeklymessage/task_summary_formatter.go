package weeklymessage

import (
	"fmt"
	"strings"
	"time"

	"kvantized-bot/internal/domain/taskboard"
)

type TaskSummaryFormatter struct {
	Now func() time.Time
}

func (f TaskSummaryFormatter) Format(summary taskboard.TaskSummary) string {
	year, week := f.now().ISOWeek()
	title := fmt.Sprintf("📅 Week %d, year %d", week, year)

	if !summary.HasActiveTasks() {
		return title + "\n\nNo active tasks on board"
	}

	var message strings.Builder
	message.WriteString(title)
	message.WriteString("\n\nTrello tasks by column:\n\n")

	for _, column := range summary.ColumnCounts() {
		fmt.Fprintf(&message, "- %s: %d\n", column.ColumnName, column.TaskCount)
	}

	return strings.TrimRight(message.String(), "\n")
}

func (f TaskSummaryFormatter) now() time.Time {
	if f.Now != nil {
		return f.Now()
	}

	return time.Now()
}
