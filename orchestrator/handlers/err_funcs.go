package handlers

import "net/http"

func HTTPErrorUnavailable(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusServiceUnavailable)+"\n"+err.Error(), http.StatusServiceUnavailable)
}

func HTTPOk(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}

func SendJSON(w http.ResponseWriter, jsonObj any) {
	if err := Render.JSON(w, http.StatusOK, jsonObj); err != nil {
		HTTPErrorUnavailable(w, err)
	}
}
