package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Frozelo/startupFeed/internal/models"
	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
)

type ProjectService interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id int64) (*models.Project, error)
	SetLike(ctx context.Context, projectId int64) error
}

type Handlers struct {
	projectService ProjectService
}

func New(projectService ProjectService) *Handlers {
	return &Handlers{projectService: projectService}
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var newProject *models.Project
	if err := json.NewDecoder(r.Body).Decode(&newProject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := h.projectService.Create(r.Context(), newProject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) FindById(w http.ResponseWriter, r *http.Request) {
	var headers map[string]string
	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		httpwriter.ErrorResponse(w, http.StatusBadRequest, err, headers)
	}
	project, err := h.projectService.FindByID(r.Context(), projectId)
	if err != nil {
		httpwriter.ErrorResponse(
			w,
			http.StatusInternalServerError,
			err,
			headers,
		)
	}
	httpwriter.SuccessResponse(w, http.StatusOK, project, headers)
}

func (h *Handlers) SetLike(w http.ResponseWriter, r *http.Request) {
	projectIdStr := chi.URLParam(r, "projectId")
	projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := h.projectService.SetLike(r.Context(), projectId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
