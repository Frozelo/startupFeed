package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Frozelo/startupFeed/internal/dto"
	"github.com/Frozelo/startupFeed/internal/middlewares"
	"github.com/Frozelo/startupFeed/internal/models"
	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
	"github.com/Frozelo/startupFeed/pkg/jwt"
)

type ProjectService interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id int64) (*models.Project, error)
	SetLike(ctx context.Context, projectId int64) error
	SetDescription(
		ctx context.Context,
		projectId int64,
		userId int64,
		updateProjectDto *dto.UpdateProjectDTO,
	) error
	CreateFeedback(
		ctx context.Context,
		projectId int64,
		feedback *dto.CreateFeedbackDTO,
	) error
}

type UserService interface {
	Register(ctx context.Context, userDTO *dto.CreateUserDTO) error
	Login(ctx context.Context, loginDTO *dto.LoginUserDTO) (*models.User, error)
}

type Handlers struct {
	projectService ProjectService
	userService    UserService
}

func New(userService UserService, projectService ProjectService) *Handlers {
	return &Handlers{userService: userService, projectService: projectService}
}

// Обработчик создания проекта
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var newProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&newProject); err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request payload",
			nil,
		)
		return
	}

	if err := h.projectService.Create(r.Context(), &newProject); err != nil {
		httpwriter.Error(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to create project",
			nil,
		)
		return

	}

	httpwriter.Success(
		w,
		http.StatusCreated,
		"Project created successfully",
		nil,
	)
}

// Обработчик поиска проекта по ID
func (h *Handlers) FindById(w http.ResponseWriter, r *http.Request) {
	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid project ID",
			nil,
		)
		return
	}

	project, err := h.projectService.FindByID(r.Context(), projectId)
	if err != nil {
		httpwriter.Error(w, http.StatusNotFound, err, "Project not found", nil)
		return
	}

	httpwriter.Success(w, http.StatusOK, project, nil)
}

// Обработчик для лайков проекта
func (h *Handlers) SetLike(w http.ResponseWriter, r *http.Request) {
	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid project ID",
			nil,
		)
		return
	}

	if err := h.projectService.SetLike(r.Context(), projectId); err != nil {
		httpwriter.Error(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to like project",
			nil,
		)
		return
	}

	httpwriter.Success(w, http.StatusOK, "Project liked successfully", nil)
}

func (h *Handlers) SetDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updateProjectDto dto.UpdateProjectDTO
	if err := json.NewDecoder(r.Body).Decode(&updateProjectDto); err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request payload",
			nil,
		)
		return
	}

	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid project ID",
			nil,
		)
		return
	}

	userId, ok := middlewares.GetUserIDFromContext(ctx)
	log.Printf("the user id is %v", userId)
	if !ok {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid user ID",
			nil,
		)
		return
	}

	if err := h.projectService.SetDescription(ctx, projectId, userId, &updateProjectDto); err != nil {
		httpwriter.Error(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to update project description",
			nil,
		)
		return
	}

	httpwriter.Success(
		w,
		http.StatusOK,
		"Project description updated successfully",
		nil,
	)
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var newProject *dto.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&newProject); err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request payload",
			nil,
		)
		return
	}
	if err := h.userService.Register(r.Context(), newProject); err != nil {
		httpwriter.Error(w, http.StatusNotFound, err, err.Error(), nil)
		return
	}

	httpwriter.Success(w, http.StatusOK, "User registered successfully", nil)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var logedUser dto.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&logedUser); err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request payload",
			nil,
		)
		return
	}

	newUser, err := h.userService.Login(r.Context(), &logedUser)
	if err != nil {
		if err.Error() == "user not found" {
			httpwriter.Error(w, http.StatusNotFound, err, err.Error(), nil)
		} else {
			httpwriter.Error(w, http.StatusUnauthorized, err, err.Error(), nil)
		}
		return
	}

	token, err := jwt.CreateToken(newUser)
	if err != nil {
		httpwriter.Error(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to generate token",
			nil,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwtToken",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	httpwriter.Success(w, http.StatusOK, token, nil)
}

func (h *Handlers) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newFeedback dto.CreateFeedbackDTO
	if err := json.NewDecoder(r.Body).Decode(&newFeedback); err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid request payload",
			nil,
		)
		return
	}

	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		httpwriter.Error(
			w,
			http.StatusBadRequest,
			err,
			"Invalid project ID",
			nil,
		)
		return
	}

	if err := h.projectService.CreateFeedback(ctx, projectId, &newFeedback); err != nil {
		httpwriter.Error(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to create feedback",
			nil,
		)
		return

	}

	httpwriter.Success(
		w,
		http.StatusCreated,
		"Feedback created successfully",
		nil,
	)
}
