package services

import (
	"context"
	"reviews/internal/models"
	"reviews/internal/repositories"
)

type ReviewService struct {
	Repo *repositories.ReviewRepository
}

func (s *ReviewService) Create(ctx context.Context, rev models.Reviews) (models.Reviews, error) {
	return s.Repo.CreateReview(ctx, rev)
}

func (s *ReviewService) GetAll(ctx context.Context) ([]models.Reviews, error) {
	return s.Repo.GetReviews(ctx)
}

func (s *ReviewService) GetByID(ctx context.Context, id int) (models.Reviews, error) {
	return s.Repo.GetReviewByID(ctx, id)
}

func (s *ReviewService) Update(ctx context.Context, rev models.Reviews) (models.Reviews, error) {
	return s.Repo.UpdateReview(ctx, rev)
}

func (s *ReviewService) Delete(ctx context.Context, id int) error {
	return s.Repo.DeleteReview(ctx, id)
}
