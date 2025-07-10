package repository

import "reviews/internal/models"

type ReviewRepository interface {
	Create(models.Review) (models.Review, error)
	GetAll() ([]models.Review, error)
	GetByID(int) (models.Review, bool, error)
	Update(models.Review) (models.Review, bool, error)
	Delete(int) bool
}
