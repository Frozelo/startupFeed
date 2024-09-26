package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Frozelo/startupFeed/internal/handlers"
	"github.com/Frozelo/startupFeed/internal/middlewares"
	"github.com/Frozelo/startupFeed/internal/repo"
	"github.com/Frozelo/startupFeed/internal/services"
	"github.com/Frozelo/startupFeed/internal/storage"
	"github.com/Frozelo/startupFeed/pkg/logger"
)

func main() {
	// TODO Implement config logic
	l, err := logger.New(slog.LevelInfo, "file.log")
	if err != nil {
		log.Fatal("logger initialization failed")
	}

	db, err := storage.New(
		"postgresql://localhost:5432/startupfeed",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	l.Info("postgres initialized")

	redisClient := storage.NewRedisClient()
	l.Info("redis initialized", nil)
	projectRepo := repo.NewProjectRepo(db.Conn)
	userRepo := repo.NewUserRepo(db.Conn)

	projectService := services.NewProjectService(
		projectRepo,
		userRepo,
		redisClient,
	)
	userService := services.NewUserSerice(userRepo)
	handler := handlers.New(userService, projectService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CORS)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/v1", func(r chi.Router) {
		r.Post("/users/register", handler.Register)
		r.Post("/users/login", handler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.JwtAuth)

			r.Get("/projects/{projectId}", handler.FindById)
			r.Post("/projects", handler.Create)
			r.Put("/projects/{projectId}", handler.SetLike)
			r.Put("/projects/{projectId}/update", handler.SetDescription)
		})
	})

	r.Mount("/api", apiRouter)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
