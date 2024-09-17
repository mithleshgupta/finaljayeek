package persistence

import (
	"fmt"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// LocationRepository implements the repository.LocationRepository interface
type LocationRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
} 

// NewLocationRepository creates a new instance of the LocationRepository
func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) CountLocations() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.Location{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LocationRepository) GetAllLocations(page int, perPage int) ([]entity.Location, error) {
	var locations []entity.Location
	if err := r.db.Debug().Model(&entity.Location{}).Limit(perPage).Offset((page - 1) * perPage).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// GetLocationByID retrieves a location by its ID
func (r *LocationRepository) GetLocationByID(id uint64) (*entity.Location, error) {
	// Location struct to store the retrieved location data
	var location entity.Location
	// Find the location by its ID and store the data in the location struct
	if err := r.db.Debug().Where("id = ?", id).Take(&location).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the location data and nil error
	return &location, nil
}

// GetLocationByCoordinates retrieves a location by its coordinates
func (r *LocationRepository) GetLocationByCoordinates(longitude, latitude float64, radius *float64) (*entity.Location, error) {
	// Location struct to store the retrieved location data
	var location entity.Location

	// Prepare the base query
	query := r.db.Debug().Table("locations").Select("*")

	// If the radius is provided, add the ST_DWithin condition
	if radius != nil {
		query = query.Where(
			"ST_DWithin(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(?, ?)::geography, ?)",
			longitude, latitude, *radius,
		)
	}

	// Complete the query by ordering and taking the location
	query = query.Order(
		fmt.Sprintf("ST_Distance(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(%.6f, %.6f)::geography)", longitude, latitude),
	).Take(&location)

	// Execute the query and store the data in the location struct
	if err := query.Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}

	// Return the location data and nil error
	return &location, nil
}
