package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// RatingRepository implements repository.RatingRepository
// and handles CRUD operations for Rating entities
type RatingRepository struct {
	db *gorm.DB
}

// NewRatingRepository returns a new instance of RatingRepository
func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db}
}

// CreateRating creates a new rating in the database
func (r *RatingRepository) CreateRating(rating *entity.Rating) (*entity.Rating, error) {
	if err := r.db.Debug().Model(&rating).Create(&rating).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&rating).Take(&rating).Error; err != nil {
		return nil, err
	}

	return rating, nil
}
