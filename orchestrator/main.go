package main

import (
	"fmt"
	"net/http"
	"os"

	"consts"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port, ok := os.LookupEnv(consts.EnvPort)
	if !ok {
		port = consts.OrchestratorDefaultPort
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

	})

	err := http.ListenAndServe(fmt.Sprintf(`:%s`, port), r)
	if err != nil {
		panic(err)
	}
}
