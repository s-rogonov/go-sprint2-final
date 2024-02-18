package handlers

import (
	"fmt"

	"dbprovider/models"
)

func query2map(q *models.Query) map[string]any {
	if q.HasError {
		return map[string]any{
			"id":       q.ID,
			"expr":     q.Expression,
			"status":   "error",
			"errorMsg": q.BadMessage,
		}
	} else if q.IsDone {
		return map[string]any{
			"id":     q.ID,
			"expr":   q.Expression,
			"status": "finished",
			"result": q.Result,
		}
	} else {
		total := len(q.Tasks)
		finished := 0
		for _, t := range q.Tasks {
			if t.IsDone {
				finished += 1
			}
		}
		total -= int(q.PlainNumbers)
		finished -= int(q.PlainNumbers)
		return map[string]any{
			"id":       q.ID,
			"expr":     q.Expression,
			"status":   "in-progress",
			"progress": fmt.Sprintf(`%v/%v`, finished, total),
		}
	}
}
