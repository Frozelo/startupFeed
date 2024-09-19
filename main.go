package main

import (
	"github.com/Frozelo/startupFeed/internal/handlers"
	"github.com/Frozelo/startupFeed/internal/repo"
	"github.com/Frozelo/startupFeed/internal/services"
	"github.com/Frozelo/startupFeed/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	db, err := storage.New("postgresql://localhost:5432/startupfeed")
	if err != nil {
		log.Fatal(err)
	}
	projectRepo := repo.NewProjectRepo(db.Conn)
	projectService := services.NewProjectService(projectRepo)
	handler := handlers.New(projectService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/projects/{projectId}", handler.FindById)
	r.Post("/projects", handler.Create)
	r.Post("/projects/{projectId}", handler.SetLike)

	http.ListenAndServe(":8080", r)
}
