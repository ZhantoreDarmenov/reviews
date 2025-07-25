package repositories

import (
	"context"
	"database/sql"
	"errors"
	"reviews/internal/models"

	_ "strings"

	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) ClearSession(ctx context.Context, id string) error {
	query := `
                UPDATE users
                SET refresh_token = NULL, expires_at = NULL
                WHERE id = ?
        `
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (r *UserRepository) SetSession(ctx context.Context, id string, session models.Session) error {

	query := `
		UPDATE users 
		SET refresh_token = ?, expires_at = ? 
		WHERE id = ?
	`

	result, err := r.DB.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (r *UserRepository) GetSession(ctx context.Context, id string) (models.Session, error) {
	query := `
		SELECT refresh_token, expires_at
		FROM users
		WHERE id = ?
	`

	var session models.Session
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, errors.New("no session found for the user")
		}
		return session, err
	}

	return session, nil
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	query := `
        SELECT id, name, password, role, created_at, updated_at
        FROM users
        WHERE name = ?
    `
	err := r.DB.QueryRowContext(ctx, query, login).Scan(
		&user.ID, &user.Name, &user.Password,
		&user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// CreateUser inserts a new user record and returns the created user with ID set.
func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
        INSERT INTO users (name, password, role, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `
	user.CreatedAt = time.Now()
	user.UpdatedAt = &user.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, user.Name, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	user.ID = int(id)
	return user, nil

}
