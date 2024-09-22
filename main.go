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
		"postgresql://localhost:5432/startupfeed",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := storage.NewRedisClient()
	log.Println("redis successfully initialized", redisClient)

	projectRepo := repo.NewProjectRepo(db.Conn)
	projectService := services.NewProjectService(projectRepo, redisClient)

	userRepo := repo.NewUserRepo(db.Conn)
	userService := services.NewUserSerice(userRepo)
	handler := handlers.New(userService, projectService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiRouter := chi.NewRouter()
	apiRouter.Route("/v1", func(r chi.Router) {
		r.Get("/projects/{projectId}", handler.FindById)
		r.Post("/projects", handler.Create)
		r.Put("/projects/{projectId}", handler.SetLike)
		r.Put("/projects/{projectId}/update", handler.SetDescription)

		r.Post("/users/register", handler.Register)
		r.Post("/users/login", handler.Login)
	})

	r.Mount("/api", apiRouter)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
