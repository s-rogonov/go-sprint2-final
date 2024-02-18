package handlers

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"dbprovider"
	"dbprovider/models"
	"parser"
)

func PostQuery(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	expr := string(bytes)
	durations, err := AccessDurations()
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	q, err := parser.ParseExpression(expr, durations)
	if err != nil {
		q = &models.Query{
			Expression: expr,
			BadMessage: err.Error(),
		}
	} else {
		q.Expression = expr
	}

	if err := dbprovider.Manager.NewQuery(q); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	jsonObj := map[string]any{
		"id":       q.ID,
		"hasError": q.HasError,
		"errorMsg": q.BadMessage,
	}

	SendJSON(w, jsonObj)
}

type OperationData struct {
	Id   uint          `json:"id"`
	Op   string        `json:"op"`
	Time time.Duration `json:"time"`
	Args []float64     `json:"args"`
}

func PostTasks(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	amount, err := strconv.ParseUint(string(bytes), 10, 0)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	durations, err := AccessDurations()
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	workers, err := dbprovider.Manager.CreateWorkers(uint(amount))
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	data := make([]OperationData, len(workers))
	for i, wk := range workers {
		args := make([]float64, len(wk.TargetTask.Subtasks))
		for _, st := range wk.TargetTask.Subtasks {
			args[st.Index] = st.Result
		}

		data[i] = OperationData{
			Id:   wk.ID,
			Op:   wk.TargetTask.Operation,
			Time: durations[wk.TargetTask.Operation],
			Args: args,
		}
	}

	SendJSON(w, data)
}
