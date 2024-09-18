package main

import (
	"github.com/Frozelo/startupFeed/internal/handlers"
	"github.com/Frozelo/startupFeed/internal/models"
	"github.com/Frozelo/startupFeed/internal/repo"
	"github.com/Frozelo/startupFeed/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	users := map[string]*models.User{"test": {
		Id:       1,
		Username: "test",
	}}

	memoryRepo := repo.NewInMemoryRepository(users)
	userService := services.NewUserService(memoryRepo)

	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/comments", userHandler.AddComment)

	http.ListenAndServe(":8080", r)
}
