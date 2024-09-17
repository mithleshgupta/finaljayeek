package repository

import "github.com/OmarBader7/web-service-jayeek/domain/entity"

// RatingRepository defines the methods that a rating repository should implement
type RatingRepository interface {
	CreateRating(*entity.Rating) (*entity.Rating, error)
}
