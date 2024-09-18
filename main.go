package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Frozelo/startupFeed/internal/handlers"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", handlers.HealthCheckHandler)

	http.ListenAndServe(":8080", r)
}
