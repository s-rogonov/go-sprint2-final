package main

import (
	"fmt"
	"net/http"
	"os"

	"consts"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"handlers"
)

func main() {
	port, ok := os.LookupEnv(consts.EnvPort)
	if !ok {
		port = consts.OrchestratorDefaultPort
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Put("/query", handlers.PutQuery)
	r.Put("/timings", handlers.PutTimings)
	r.Put("/result", handlers.PutResult)

	r.Put("/query", handlers.PostQuery)
	r.Put("/tasks", handlers.PostTasks)

	err := http.ListenAndServe(fmt.Sprintf(`:%s`, port), r)
	if err != nil {
		panic(err)
	}
}
