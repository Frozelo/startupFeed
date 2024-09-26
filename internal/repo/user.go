package repo

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/Frozelo/startupFeed/internal/models"
)

type UserRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	user := &models.User{}

	query := `SELECT * FROM users WHERE email = $1`
	if err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.RegistrationDate,
		&user.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}

	return user, nil
}

func (r *UserRepo) Create(
	ctx context.Context,
	user *models.User,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Println(err)
			tx.Rollback(ctx)
		}
	}()

	query := `INSERT INTO users(name, email, password_hash, role) VALUES($1, $2, $3, $4)`
	_, err = tx.Exec(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetAuthors(
	ctx context.Context,
	projectId int64,
) ([]int64, error) {
	query := `SELECT user_id FROM project_authors WHERE project_id = $1`
	rows, err := r.db.Query(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authorsId []int64

	for rows.Next() {
		var authorId int64

		if err := rows.Scan(&authorId); err != nil { // Передаем указатель на authorId
			return nil, err
		}
		authorsId = append(authorsId, authorId) // Добавляем authorId в срез
	}

	if err := rows.Err(); err != nil { // Проверяем на наличие ошибок после итерации
		return nil, err
	}

	return authorsId, nil
}
