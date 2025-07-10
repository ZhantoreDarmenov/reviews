package memory

import (
	"reviews/internal/models"
	"sync"
)

type MemoryReviewRepository struct {
	mu      sync.Mutex
	nextID  int
	reviews []models.Review
}

func New() *MemoryReviewRepository {
	return &MemoryReviewRepository{nextID: 1}
}

func (r *MemoryReviewRepository) Create(rev models.Review) (models.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rev.ID = r.nextID
	r.nextID++
	r.reviews = append(r.reviews, rev)
	return rev, nil
}

func (r *MemoryReviewRepository) GetAll() ([]models.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	res := make([]models.Review, len(r.reviews))
	copy(res, r.reviews)
	return res, nil
}

func (r *MemoryReviewRepository) GetByID(id int) (models.Review, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rev := range r.reviews {
		if rev.ID == id {
			return rev, true, nil
		}
	}
	return models.Review{}, false, nil
}

func (r *MemoryReviewRepository) Update(upd models.Review) (models.Review, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rev := range r.reviews {
		if rev.ID == upd.ID {
			r.reviews[i] = upd
			return upd, true, nil
		}
	}
	return models.Review{}, false, nil
}

func (r *MemoryReviewRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, rev := range r.reviews {
		if rev.ID == id {
			r.reviews = append(r.reviews[:i], r.reviews[i+1:]...)
			return true
		}
	}
	return false
}
