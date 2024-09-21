package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Frozelo/startupFeed/internal/dto"
	"github.com/Frozelo/startupFeed/internal/models"
	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
)

type ProjectService interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id int64) (*models.Project, error)
	SetLike(ctx context.Context, projectId int64) error
	SetDescription(
		ctx context.Context,
		projectId int64,
		updateProjectDto *dto.UpdateProjectDTO,
	) error
}

type Handlers struct {
	projectService ProjectService
}

func New(projectService ProjectService) *Handlers {
	return &Handlers{projectService: projectService}
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

// Обработчик обновления описания проекта
func (h *Handlers) SetDescription(w http.ResponseWriter, r *http.Request) {
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

	if err := h.projectService.SetDescription(r.Context(), projectId, &updateProjectDto); err != nil {
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
