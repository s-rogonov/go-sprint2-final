package main

import (
	"fmt"
	"net/http"
	"os"

	"consts"
	"dbprovider"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"handlers"
)

func main() {
	dbprovider.InitConnection()

	port, ok := os.LookupEnv(consts.EnvPort)
	if !ok {
		port = consts.OrchestratorDefaultPort
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Put("/query", handlers.PutQuery)
	r.Put("/timings", handlers.PutTimings)
	r.Put("/result", handlers.PutResult)

	r.Post("/query", handlers.PostQuery)
	r.Post("/tasks", handlers.PostTasks)

	r.Get("/timings", handlers.GetTimings)

	r.Route("/queries", func(r chi.Router) {
		r.Get("/", handlers.GetLastQueries)
		r.Get("/{id}", handlers.GetQuery)
	})

	r.Route("/workers", func(r chi.Router) {
		r.Get("/", handlers.GetLastWorkers)
		r.Get("/{id}", handlers.GetWorker)
	})

	err := http.ListenAndServe(fmt.Sprintf(`:%s`, port), r)
	if err != nil {
		panic(err)
	}
}
