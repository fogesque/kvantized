package taskboard

type ColumnID string

type Column struct {
	ID   ColumnID
	Name string
}

type Card struct {
	ColumnID ColumnID
}

type ColumnTaskCount struct {
	ColumnName string
	TaskCount  int
}

type TaskSummary struct {
	columns []ColumnTaskCount
}

func NewTaskSummary(columns []Column, cards []Card) TaskSummary {
	cardCounts := make(map[ColumnID]int, len(columns))
	for _, card := range cards {
		cardCounts[card.ColumnID]++
	}

	columnCounts := make([]ColumnTaskCount, 0, len(columns))
	for _, column := range columns {
		count := cardCounts[column.ID]
		if count == 0 {
			continue
		}

		columnCounts = append(columnCounts, ColumnTaskCount{
			ColumnName: column.Name,
			TaskCount:  count,
		})
	}

	return TaskSummary{columns: columnCounts}
}

func (s TaskSummary) ColumnCounts() []ColumnTaskCount {
	return append([]ColumnTaskCount(nil), s.columns...)
}

func (s TaskSummary) HasActiveTasks() bool {
	return len(s.columns) > 0
}
