package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"dbprovider"
	"dbprovider/models"
	"parser"
)

func PutTimings(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	var td TimingsData
	if err := json.Unmarshal(bytes, &td); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	tm := &models.Timings{}
	if err := dbprovider.Manager.GetTimings(tm); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	tm.Factor = td.Factor
	tm.Addition = time.Duration(float32(1*time.Second) * td.Add)
	tm.Subtraction = time.Duration(float32(1*time.Second) * td.Sub)
	tm.Multiplication = time.Duration(float32(1*time.Second) * td.Mul)
	tm.Division = time.Duration(float32(1*time.Second) * td.Div)

	if err := dbprovider.Manager.UpdateTimings(tm); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	HTTPOk(w)
}

type Result struct {
	Id     uint    `json:"id"`
	Result float64 `json:"result"`
}

func PutResult(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	var res Result
	if err := json.Unmarshal(bytes, &res); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	if err := dbprovider.Manager.SetWorkResult(res.Id, res.Result); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	HTTPOk(w)
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
