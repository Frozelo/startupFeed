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
