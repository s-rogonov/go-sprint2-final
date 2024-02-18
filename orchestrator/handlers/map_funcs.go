package handlers

import (
	"fmt"
	"log"
	"time"

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

func worker2map(w *models.Worker, factor float32) map[string]any {
	expr := fmt.Sprintf(
		`%v %v %v`,
		w.TargetTask.Subtasks[0].Result,
		w.TargetTask.Operation,
		w.TargetTask.Subtasks[1].Result)

	eta := w.CreatedAt.Add(w.TargetTask.Duration)
	deadline := w.CreatedAt.Add(time.Duration(factor * float32(w.TargetTask.Duration)))
	now := time.Now()

	log.Println(now, eta, deadline, w.TargetTask.Duration)

	if w.IsDone {
		return map[string]any{
			"id":     w.ID,
			"expr":   expr,
			"status": "finished",
			"result": w.Result,
		}
	} else if now.Before(eta) {
		left := float32(eta.Sub(now)) / float32(1*time.Second)
		return map[string]any{
			"id":     w.ID,
			"expr":   expr,
			"status": "computing",
			"left":   left,
		}
	} else if now.Before(deadline) {
		left := float32(deadline.Sub(now)) / float32(1*time.Second)
		return map[string]any{
			"id":     w.ID,
			"expr":   expr,
			"status": "retrieving",
			"left":   left,
		}
	} else {
		return map[string]any{
			"id":       w.ID,
			"expr":     expr,
			"status":   "timeout",
			"deadline": deadline,
		}
	}
}
