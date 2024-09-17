package entity

// Rating represents a rating in the system
type Rating struct {
	ID                uint64   `gorm:"primary_key;auto_increment" json:"id"`
	UserID            uint64   `json:"user_id" validate:"required,numeric"`
	OrderID           uint64   `json:"order_id" validate:"required,numeric"`
	FastRating        *float64 `gorm:"default:null" json:"fast_rating" validate:"omitempty,numeric"`
	ExperienceRating  *float64 `gorm:"default:null" json:"experience_rating" validate:"omitempty,numeric"`
	RecommendedRating *float64 `gorm:"default:null" json:"recommended_rating" validate:"omitempty,numeric"`
}

type RatingPublicData struct {
	ID                uint64   `json:"id"`
	UserID            uint64   `json:"user_id"`
	OrderID           uint64   `json:"order_id"`
	FastRating        *float64 `json:"fast_rating"`
	ExperienceRating  *float64 `json:"experience_rating"`
	RecommendedRating *float64 `json:"recommended_rating"`
}

// PublicData returns a copy of the rating's public information
func (r *Rating) PublicData() interface{} {
	return &RatingPublicData{
		ID:                r.ID,
		UserID:            r.UserID,
		OrderID:           r.OrderID,
		FastRating:        r.FastRating,
		ExperienceRating:  r.ExperienceRating,
		RecommendedRating: r.RecommendedRating,
	}
}
