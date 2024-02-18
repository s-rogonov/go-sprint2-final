package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"dbprovider"
	"dbprovider/models"
	"parser"
)

func PutTimings(writer http.ResponseWriter, request *http.Request) {

}

func PutResult(writer http.ResponseWriter, request *http.Request) {

}

type QueryFix struct {
	Id   uint   `json:"id"`
	Expr string `json:"expr"`
}

func PutQuery(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	var qfix QueryFix
	if err := json.Unmarshal(bytes, &qfix); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	durations, err := AccessDurations()
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	q, err := parser.ParseExpression(qfix.Expr, durations)
	if err != nil {
		q = &models.Query{
			Expression: qfix.Expr,
			BadMessage: err.Error(),
		}
	} else {
		q.Expression = qfix.Expr
	}
	q.ID = qfix.Id

	if err := dbprovider.Manager.UpdateQuery(q); err != nil {
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
