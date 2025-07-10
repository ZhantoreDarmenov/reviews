package service

import (
	"reviews/internal/models"
	"reviews/internal/repository"
)

type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(r repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: r}
}

func (s *ReviewService) Create(rev models.Review) (models.Review, error) {
	return s.repo.Create(rev)
}

func (s *ReviewService) GetAll() ([]models.Review, error) {
	return s.repo.GetAll()
}

func (s *ReviewService) GetByID(id int) (models.Review, bool, error) {
	return s.repo.GetByID(id)
}

func (s *ReviewService) Update(rev models.Review) (models.Review, bool, error) {
	return s.repo.Update(rev)
}

func (s *ReviewService) Delete(id int) bool {
	return s.repo.Delete(id)
}
