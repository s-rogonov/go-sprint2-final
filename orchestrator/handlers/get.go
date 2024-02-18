package handlers

import (
	"net/http"
	"time"

	"dbprovider"
	"dbprovider/models"
)

// func GetQueries(writer http.ResponseWriter, request *http.Request) {
//
// }
//

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
