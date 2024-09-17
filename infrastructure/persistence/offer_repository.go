package persistence

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// OfferRepository implements the repository.OfferRepository interface
type OfferRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewOfferRepository creates a new instance of the OfferRepository
func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

// CreateOffer creates a new offer in the database
func (r *OfferRepository) CreateOffer(offer *entity.Offer) (*entity.Offer, error) {
	if err := r.db.Debug().Model(&offer).Create(&offer).Error; err != nil {
		return nil, err
	}
	if err := r.db.Debug().Model(&offer).Take(&offer).Error; err != nil {
		return nil, err
	}

	return offer, nil
}

// UpdateOffer updates the offer
func (r *OfferRepository) UpdateOfferByID(id uint64, offer *entity.Offer) (*entity.Offer, error) {
	if err := r.db.Debug().Model(&offer).Updates(offer).Where("id = ?", id).Error; err != nil {
		return nil, err
	}

	return offer, nil
}

func (r *OfferRepository) CountOffersByDriverIDAndStatus(driverID uint64, status entity.OfferStatus) (int64, error) {
	var count int64
	if err := r.db.Debug().Table("offers").Where("driver_id = ?", driverID).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OfferRepository) CountOffersByStatusAndUserID(status entity.OfferStatus, userID uint64) (int64, error) {
	var count int64
	if err := r.db.Debug().Table("offers").Joins("left join orders on offers.order_id = orders.id").Where("offers.status = ?", status).Where("orders.user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OfferRepository) GetAllOffersByStatusAndUserID(status entity.OfferStatus, userID uint64, page int, perPage int) ([]entity.Offer, error) {
	var offers []entity.Offer
	if err := r.db.Debug().Table("offers").Joins("left join orders on offers.order_id = orders.id").Where("offers.status = ?", status).Where("orders.user_id = ?", userID).Limit(perPage).Offset((page - 1) * perPage).Order("offers.created_at desc").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Order").Preload("Order.Location").Preload("Order.User").Preload("Order.User.Location").Preload("Order.Driver").Preload("Order.Driver.User").Preload("Order.Driver.User.Location").Preload("Order.Driver.TransportationMode").Preload("Order.Category").Preload("Order.Size").Preload("Order.DeliveryTime").Preload("Order.ShipmentContents").Preload("Order.ExtraServices").Preload("Order.Destination").Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}

func (r *OfferRepository) GetOfferByIDAndUserID(id uint64, userID uint64) (*entity.Offer, error) {
	// Offer struct to store the retrieved offer data
	var offer entity.Offer
	// Find the offer by its ID and store the data in the offer struct
	if err := r.db.Debug().Table("offers").Joins("left join orders on offers.order_id = orders.id").Where("offers.id = ?", id).Where("orders.user_id = ?", userID).Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Order").Preload("Order.Location").Preload("Order.User").Preload("Order.User.Location").Preload("Order.Driver").Preload("Order.Driver.User").Preload("Order.Driver.User.Location").Preload("Order.Driver.TransportationMode").Preload("Order.Category").Preload("Order.Size").Preload("Order.DeliveryTime").Preload("Order.ShipmentContents").Preload("Order.ExtraServices").Preload("Order.Destination").Take(&offer).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the offer data and nil error
	return &offer, nil
}

func (r *OfferRepository) GetAllOffersByStatusAndOrderID(status entity.OfferStatus, orderID uint64) ([]entity.Offer, error) {
	var offers []entity.Offer
	if err := r.db.Debug().Where("status = ?", status).Where("order_id = ?", orderID).Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Order").Preload("Order.Location").Preload("Order.User").Preload("Order.User.Location").Preload("Order.Driver").Preload("Order.Driver.User").Preload("Order.Driver.User.Location").Preload("Order.Driver.TransportationMode").Preload("Order.Category").Preload("Order.Size").Preload("Order.DeliveryTime").Preload("Order.ShipmentContents").Preload("Order.ExtraServices").Preload("Order.Destination").Find(&offers).Error; err != nil {
		return nil, err
	}
	return offers, nil
}
