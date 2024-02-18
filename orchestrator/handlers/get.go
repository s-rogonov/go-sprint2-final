package handlers

import (
	"net/http"
	"strconv"
	"time"

	"dbprovider"
	"dbprovider/models"
	"github.com/go-chi/chi/v5"
)

func GetLastWorkers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "3"
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	ws, err := dbprovider.Manager.GetWorkers(uint(limit))
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	tm := &models.Timings{}
	if err := dbprovider.Manager.GetTimings(tm); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	var answer []map[string]any
	for _, wk := range ws {
		answer = append(answer, worker2map(wk, tm.Factor))
	}

	SendJSON(w, answer)
}

func GetWorker(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	wk, err := dbprovider.Manager.GetWorker(uint(id))
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	tm := &models.Timings{}
	if err := dbprovider.Manager.GetTimings(tm); err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	SendJSON(w, worker2map(wk, tm.Factor))
}

func GetLastQueries(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "3"
	}

	limit, err := strconv.ParseUint(limitStr, 10, 0)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	qs, err := dbprovider.Manager.GetQueries(uint(limit))
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	var answer []map[string]any
	for _, q := range qs {
		answer = append(answer, query2map(q))
	}

	SendJSON(w, answer)
}

func GetQuery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	q, err := dbprovider.Manager.GetQuery(uint(id))
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	SendJSON(w, query2map(q))
}

type TimingsData struct {
	Factor float32 `json:"factor"`
	Add    float32 `json:"add"`
	Sub    float32 `json:"sub"`
	Mul    float32 `json:"mul"`
	Div    float32 `json:"div"`
}

func GetTimings(w http.ResponseWriter, _ *http.Request) {
	tm := &models.Timings{}
	err := dbprovider.Manager.GetTimings(tm)
	if err != nil {
		HTTPErrorUnavailable(w, err)
		return
	}

	td := TimingsData{
		Factor: tm.Factor,
		Add:    float32(tm.Addition) / float32(1*time.Second),
		Sub:    float32(tm.Subtraction) / float32(1*time.Second),
		Mul:    float32(tm.Multiplication) / float32(1*time.Second),
		Div:    float32(tm.Division) / float32(1*time.Second),
	}

	SendJSON(w, td)
}

//
// func HelloWorld(w http.ResponseWriter, r *http.Request) {
// 	_, err := w.Write([]byte("welcome"))
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
// 		return
// 	}
//
// }
