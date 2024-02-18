package handlers

import (
	"io"
	"net/http"

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

func PostTasks(writer http.ResponseWriter, request *http.Request) {

}
