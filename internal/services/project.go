package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/Frozelo/startupFeed/internal/dto"
	"github.com/Frozelo/startupFeed/internal/models"
)

type ProjectRepo interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id int64) (*models.Project, error)
	UpdateLikes(ctx context.Context, project *models.Project) error
	UpdateDescription(ctx context.Context, project *models.Project) error
	CreateFeedback(
		ctx context.Context,
		projectId int64,
		feedback *models.Feedback,
	) error
	GetFeedbacksByProjectId(
		ctx context.Context,
		projectId int64,
	) ([]*models.Feedback, error)
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value any, expiration time.Duration) error
	Delete(key string) error
}

type ProjectService struct {
	projectRepo ProjectRepo
	userRepo    UserRepo
	cache       Cache
}

func NewProjectService(
	projectRepo ProjectRepo,
	userRepo UserRepo,
	cache Cache,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
		cache:       cache,
	}
}

func (s *ProjectService) Create(
	ctx context.Context,
	project *models.Project,
) error {
	return s.projectRepo.Create(ctx, project)
}

// Find project by ID with Redis cache
func (s *ProjectService) FindByID(
	ctx context.Context,
	id int64,
) (*models.Project, error) {
	cacheKey := generateProjectCacheKey(id)

	if cachedData, err := s.cache.Get(cacheKey); err == nil {
		var project models.Project
		if err := json.Unmarshal([]byte(cachedData), &project); err == nil {
			log.Print("found it in cahche!")
			return &project, nil
		}
	}
	log.Println("not found it in cache")

	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	feedbacks, err := s.projectRepo.GetFeedbacksByProjectId(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	project.Feedbacks = feedbacks

	log.Println("found it in database!")

	if project == nil {
		return nil, errors.New("project not found")
	}

	if jsonData, err := json.Marshal(project); err == nil {
		s.cache.Set(cacheKey, jsonData, 10*time.Minute)
	}

	return project, nil
}

// Set like for a project and update cache
func (s *ProjectService) SetLike(ctx context.Context, projectId int64) error {
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

	cacheKey := generateProjectCacheKey(projectId)
	if jsonData, err := json.Marshal(project); err == nil {
		s.cache.Set(cacheKey, jsonData, 10*time.Minute)
	}
	log.Println("New cache set!")

	return nil
}

func (s *ProjectService) SetDescription(
	ctx context.Context,
	projectId int64,
	userId int64,
	updateProjectDto *dto.UpdateProjectDTO,
) error {
	project, err := s.projectRepo.FindByID(ctx, projectId)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("projects not found")
	}

	if !s.isAuthor(ctx, userId, projectId) {
		return errors.New("you do not have permission to edit this project")
	}

	project.Description = updateProjectDto.Description
	// TODO caching invalidation
	if err := s.projectRepo.UpdateDescription(ctx, project); err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) CreateFeedback(
	ctx context.Context,
	projectId int64,
	feedback *dto.CreateFeedbackDTO,
) error {
	newFeedback := &models.Feedback{
		UserId: feedback.UserId,
		Text:   feedback.Text,
	}
	if err := s.projectRepo.CreateFeedback(ctx, projectId, newFeedback); err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) DeleteProjectCache(
	ctx context.Context,
	projectId int64,
) error {
	cacheKey := generateProjectCacheKey(projectId)
	return s.cache.Delete(cacheKey)
}

func generateProjectCacheKey(id int64) string {
	return "project:" + strconv.FormatInt(id, 10)
}

func (s *ProjectService) isAuthor(
	ctx context.Context,
	userId int64,
	projectId int64,
) bool {
	authorsId, err := s.userRepo.GetAuthors(ctx, projectId)
	if err != nil {
		log.Printf("error fetching authors for project %d: %v", projectId, err)
		return false
	}

	return slices.Contains(authorsId, userId)
}
