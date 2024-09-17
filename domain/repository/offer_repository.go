package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// OfferRepository defines the methods for interacting with offer data
type OfferRepository interface {
	CreateOffer(*entity.Offer) (*entity.Offer, error)
	UpdateOfferByID(id uint64, offer *entity.Offer) (*entity.Offer, error)
	CountOffersByDriverIDAndStatus(driverID uint64, status entity.OfferStatus) (int64, error)
	CountOffersByStatusAndUserID(status entity.OfferStatus, userID uint64) (int64, error)
	GetAllOffersByStatusAndUserID(status entity.OfferStatus, userID uint64, page int, perPage int) ([]entity.Offer, error)
	GetOfferByIDAndUserID(id uint64, userID uint64) (*entity.Offer, error)
	GetAllOffersByStatusAndOrderID(status entity.OfferStatus, orderID uint64) ([]entity.Offer, error)
}
