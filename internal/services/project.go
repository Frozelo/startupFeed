package services

import (
	"context"
	"errors"
	"github.com/Frozelo/startupFeed/internal/models"
)

type ProjectRepo interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id int64) (*models.Project, error)
	UpdateLikes(ctx context.Context, project *models.Project) error
}

type ProjectServicee struct {
	projectRepo ProjectRepo
}

func NewProjectService(projectRepo ProjectRepo) *ProjectServicee {
	return &ProjectServicee{projectRepo: projectRepo}
}

func (s *ProjectServicee) Create(ctx context.Context, project *models.Project) error {
	return s.projectRepo.Create(ctx, project)
}

func (s *ProjectServicee) FindByID(ctx context.Context, id int64) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectServicee) SetLike(ctx context.Context, projectId int64) error {
	project, err := s.projectRepo.FindByID(ctx, projectId)
	if err != nil {
		return err
	}

	if project == nil {
		return errors.New("project not found")
	}
	project.Votes += 1
	if err := s.projectRepo.UpdateLikes(ctx, project); err != nil {
		return err
	}
	return nil
}
