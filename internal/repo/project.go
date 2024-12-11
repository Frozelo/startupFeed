package repo

import (
	"context"
	"errors"

	"github.com/Frozelo/startupFeed/internal/models"
	"github.com/jackc/pgx/v5"
)

type ProjectRepo struct {
	db *pgx.Conn
}

func NewProjectRepo(db *pgx.Conn) *ProjectRepo {
	return &ProjectRepo{db: db}
}

// Create - создание нового проекта
func (r *ProjectRepo) Create(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO projects (name, description, category_id, author_id, votes, create_date)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, project.Name, project.Description, project.CategoryId, project.AuthorId, project.Votes, project.CreateDate)
	return err
}

// GetAll - получение всех проектов
func (r *ProjectRepo) GetAll(ctx context.Context) ([]*models.Project, error) {
	query := `SELECT id, name, description, category_id, author_id, votes, create_date FROM projects`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CategoryId, &project.AuthorId, &project.Votes, &project.CreateDate); err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	return projects, nil
}

// FindByID - поиск проекта по ID
func (r *ProjectRepo) FindByID(ctx context.Context, id int64) (*models.Project, error) {
	query := `SELECT id, name, description, category_id, author_id, votes, create_date FROM projects WHERE id = $1`
	project := &models.Project{}
	err := r.db.QueryRow(ctx, query, id).Scan(&project.ID, &project.Name, &project.Description, &project.CategoryId, &project.AuthorId, &project.Votes, &project.CreateDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return project, err
}

// UpdateVotes - обновление количества голосов
func (r *ProjectRepo) UpdateVotes(ctx context.Context, id int64, votes int64) error {
	query := `UPDATE projects SET votes = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, votes, id)
	return err
}

// UpdateDescription - обновление описания проекта
func (r *ProjectRepo) UpdateDescription(ctx context.Context, id int64, description string) error {
	query := `UPDATE projects SET description = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, description, id)
	return err
}

// CreateFeedback - создание отзыва
func (r *ProjectRepo) CreateFeedback(ctx context.Context, feedback *models.Feedback) error {
	query := `INSERT INTO feedbacks (user_id, project_id, text, create_date) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, feedback.UserId, feedback.ProjectId, feedback.Text, feedback.CreateDate)
	return err
}

// GetFeedbacksByProjectId - получение отзывов для проекта
func (r *ProjectRepo) GetFeedbacksByProjectId(ctx context.Context, projectId int64) ([]*models.Feedback, error) {
	query := `SELECT id, user_id, text, create_date FROM feedbacks WHERE project_id = $1`
	rows, err := r.db.Query(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []*models.Feedback
	for rows.Next() {
		var feedback models.Feedback
		if err := rows.Scan(&feedback.ID, &feedback.UserId, &feedback.Text, &feedback.CreateDate); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, &feedback)
	}
	return feedbacks, nil
}

// DeleteProject - удаление проекта по ID
func (r *ProjectRepo) DeleteProject(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
