package helpers

import "dbprovider/models"

func DefaultsForQuery(query *models.Query) {
	query.PlainNumbers = 0
	query.IsDone = false
	query.HasError = query.BadMessage != ""
}

func DefaultsForTask(t *models.Task) {
	t.TotalSubtasks = uint(len(t.Subtasks))
	t.FinishedSubtasks = 0
	t.IsDone = len(t.Subtasks) == 0
	t.IsReady = false
}
