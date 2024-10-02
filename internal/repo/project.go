package repo

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/Frozelo/startupFeed/internal/models"
)

type ProjectRepo struct {
	db *pgx.Conn
}

func NewProjectRepo(db *pgx.Conn) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) Create(
	ctx context.Context,
	project *models.Project,
) (err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO projects(id, name, description) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, query, project.ID, project.Name, project.Description)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ProjectRepo) FindByID(
	ctx context.Context,
	id int64,
) (*models.Project, error) {
	query := `SELECT * FROM projects WHERE id = $1`
	project := &models.Project{}
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.Votes,
		&project.CreateDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Println("find by id err")
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepo) UpdateLikes(
	ctx context.Context,
	project *models.Project,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	query := `UPDATE projects SET votes = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, project.Votes, project.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepo) UpdateDescription(
	ctx context.Context,
	project *models.Project,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	log.Printf("got the new project, %v", project)
	query := `UPDATE projects SET description = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, project.Description, project.ID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepo) CreateFeedback(
	ctx context.Context,
	projectId int64,
	feedback *models.Feedback,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO feedbacks(user_id, project_id, text) VALUES ($1, $2, $3)`
	_, err = tx.Exec(
		ctx,
		query,
		feedback.UserId,
		projectId,
		feedback.Text,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ProjectRepo) GetFeedbacksByProjectId(
	ctx context.Context,
	projectId int64,
) ([]*models.Feedback, error) {
	query := `SELECT id, user_id, text, create_date FROM feedbacks WHERE project_id = $1`
	rows, err := r.db.Query(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []*models.Feedback
	for rows.Next() {
		feedback := &models.Feedback{}
		if err := rows.Scan(
			&feedback.ID,
			&feedback.UserId,
			&feedback.Text,
			&feedback.CreateDate,
		); err != nil {
			log.Println("feedbacks repo err is here!")
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedbacks, nil
}
