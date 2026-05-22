package taskboard

import "testing"

func TestNewTaskSummaryCountsCardsByColumnAndSkipsEmptyColumns(t *testing.T) {
	summary := NewTaskSummary(
		[]Column{
			{ID: "todo", Name: "Todo"},
			{ID: "review", Name: "Review"},
			{ID: "done", Name: "Done"},
		},
		[]Card{
			{ColumnID: "todo"},
			{ColumnID: "todo"},
			{ColumnID: "review"},
		},
	)

	got := summary.ColumnCounts()
	want := []ColumnTaskCount{
		{ColumnName: "Todo", TaskCount: 2},
		{ColumnName: "Review", TaskCount: 1},
	}

	if len(got) != len(want) {
		t.Fatalf("ColumnCounts() length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ColumnCounts()[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestNewTaskSummaryHasNoActiveTasksWhenEveryColumnIsEmpty(t *testing.T) {
	summary := NewTaskSummary(
		[]Column{{ID: "todo", Name: "Todo"}},
		nil,
	)

	if summary.HasActiveTasks() {
		t.Fatal("HasActiveTasks() = true, want false")
	}
}
