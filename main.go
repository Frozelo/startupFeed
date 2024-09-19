package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Frozelo/startupFeed/internal/handlers"
	"github.com/Frozelo/startupFeed/internal/repo"
	"github.com/Frozelo/startupFeed/internal/services"
	"github.com/Frozelo/startupFeed/internal/storage"
)

func main() {
	// TODO Implement config logic
	db, err := storage.New(
		"",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := storage.NewRedisClient()
	log.Println("reddis successfully initialized", redisClient)

	projectRepo := repo.NewProjectRepo(db.Conn)
	projectService := services.NewProjectService(projectRepo, redisClient)
	handler := handlers.New(projectService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiRouter := chi.NewRouter()
	apiRouter.Route("/projects", func(r chi.Router) {
		r.Get("/{projectId}", handler.FindById)
		r.Post("/", handler.Create)
		r.Post("/{projectId}", handler.SetLike)
	})

	r.Mount("/api", apiRouter)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
