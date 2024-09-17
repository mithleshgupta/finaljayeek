package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// RatingApplication handles the business logic for ratings
type RatingApplication struct {
	ratingRepo repository.RatingRepository
}

var _ RatingApplicationInterface = &RatingApplication{}

// RatingApplicationInterface defines the methods available for RatingApplication
type RatingApplicationInterface interface {
	CreateRating(balance *entity.Rating) (*entity.Rating, error)
}

// CreateUser creates a new user in the database
func (a *RatingApplication) CreateRating(rating *entity.Rating) (*entity.Rating, error) {
	return a.ratingRepo.CreateRating(rating)
}
