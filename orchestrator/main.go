package main

import (
	"net/http"

	h "orchestrator/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/timings", h.Timings)

	err := http.ListenAndServe("localhost:8181", mux)
	if err != nil {
		return
	}
}
