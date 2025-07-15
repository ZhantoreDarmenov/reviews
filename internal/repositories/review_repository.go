package repositories

import (
	"context"
	"database/sql"
	"reviews/internal/models"
	"time"
)

type ReviewRepository struct {
	DB *sql.DB
}

func (r *ReviewRepository) CreateReview(ctx context.Context, rev models.Reviews) (models.Reviews, error) {
	rev.CreatedAt = time.Now()
	rev.UpdatedAt = &rev.CreatedAt
	query := `INSERT INTO reviews (name, photo, pdf_file, industry, service, description, rating, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.DB.ExecContext(ctx, query, rev.Name, rev.Photo, rev.PdfFile, rev.Industry, rev.Service, rev.Description, rev.Rating, rev.CreatedAt, rev.UpdatedAt)
	if err != nil {
		return models.Reviews{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Reviews{}, err
	}
	rev.ID = int(id)
	return rev, nil
}

func (r *ReviewRepository) GetReviews(ctx context.Context) ([]models.Reviews, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, photo, pdf_file, industry, service, description, rating, created_at, updated_at FROM reviews`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.Reviews
	for rows.Next() {
		var rev models.Reviews
		err := rows.Scan(&rev.ID, &rev.Name, &rev.Photo, &rev.PdfFile, &rev.Industry, &rev.Service, &rev.Description, &rev.Rating, &rev.CreatedAt, &rev.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, rev)
	}
	return res, nil
}

func (r *ReviewRepository) GetReviewByID(ctx context.Context, id int) (models.Reviews, error) {
	var rev models.Reviews
	err := r.DB.QueryRowContext(ctx, `SELECT id, name, photo, pdf_file, industry, service, description, rating, created_at, updated_at FROM reviews WHERE id = ?`, id).Scan(
		&rev.ID, &rev.Name, &rev.Photo, &rev.PdfFile, &rev.Industry, &rev.Service, &rev.Description, &rev.Rating, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Reviews{}, sql.ErrNoRows
		}
		return models.Reviews{}, err
	}
	return rev, nil
}

func (r *ReviewRepository) UpdateReview(ctx context.Context, rev models.Reviews) (models.Reviews, error) {
	now := time.Now()
	rev.UpdatedAt = &now
	query := `UPDATE reviews SET name = ?, photo = ?, pdf_file = ?, industry = ?, service = ?, description = ?, rating = ?, updated_at = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, rev.Name, rev.Photo, rev.PdfFile, rev.Industry, rev.Service, rev.Description, rev.Rating, rev.UpdatedAt, rev.ID)
	if err != nil {
		return models.Reviews{}, err
	}
	return rev, nil
}

func (r *ReviewRepository) DeleteReview(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM reviews WHERE id = ?`, id)
	return err
}
