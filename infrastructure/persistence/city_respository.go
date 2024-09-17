package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// CityRepository implements the repository.CityRepository interface.
type CityRepository struct {
	db *gorm.DB
}

// NewCityRepository creates a new instance of the CityRepository.
func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) CountCities() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.City{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CityRepository) GetAllCities(page int, perPage int) ([]entity.City, error) {
	var cities []entity.City
	if err := r.db.Debug().Model(&entity.City{}).Limit(perPage).Offset((page - 1) * perPage).Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *CityRepository) GetCityByID(id uint64) (*entity.City, error) {
	var city entity.City
	if err := r.db.Debug().Where("id = ?", id).Take(&city).Error; err != nil {
		return nil, err
	}
	return &city, nil
}
