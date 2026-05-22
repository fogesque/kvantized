package weeklymessage

import (
	"testing"
	"time"

	"kvantized-bot/internal/domain/taskboard"
)

func TestTaskSummaryFormatterSkipsEmptyColumns(t *testing.T) {
	summary := taskboard.NewTaskSummary(
		[]taskboard.Column{
			{ID: "todo", Name: "Todo"},
			{ID: "review", Name: "Review"},
			{ID: "done", Name: "Done"},
		},
		[]taskboard.Card{
			{ColumnID: "todo"},
			{ColumnID: "todo"},
			{ColumnID: "todo"},
			{ColumnID: "todo"},
			{ColumnID: "review"},
		},
	)

	message := testFormatter().Format(summary)

	want := "📅 Week 1, year 2024\n\nTrello tasks by column:\n\n- Todo: 4\n- Review: 1"
	if message != want {
		t.Fatalf("Format() = %q, want %q", message, want)
	}
}

func TestTaskSummaryFormatterNoActiveTasks(t *testing.T) {
	summary := taskboard.NewTaskSummary(
		[]taskboard.Column{{ID: "todo", Name: "Todo"}},
		nil,
	)

	message := testFormatter().Format(summary)

	want := "📅 Week 1, year 2024\n\nNo active tasks on board"
	if message != want {
		t.Fatalf("Format() = %q, want %q", message, want)
	}
}

func testFormatter() TaskSummaryFormatter {
	return TaskSummaryFormatter{
		Now: func() time.Time {
			return time.Date(2024, time.January, 3, 0, 0, 0, 0, time.UTC)
		},
	}
}
