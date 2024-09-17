package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// OfferApplication handles the business logic for offers
type OfferApplication struct {
	offerRepo repository.OfferRepository
}

var _ OfferApplicationInterface = &OfferApplication{}

// OfferApplicationInterface defines the methods available for OfferApplication
type OfferApplicationInterface interface {
	CreateOffer(balance *entity.Offer) (*entity.Offer, error)
	UpdateOfferByID(id uint64, offer *entity.Offer) (*entity.Offer, error)
	CountOffersByDriverIDAndStatus(driverID uint64, status entity.OfferStatus) (int64, error)
	CountOffersByStatusAndUserID(status entity.OfferStatus, userID uint64) (int64, error)
	GetAllOffersByStatusAndUserID(status entity.OfferStatus, userID uint64, page int, perPage int) ([]entity.Offer, error)
	GetOfferByIDAndUserID(id uint64, userID uint64) (*entity.Offer, error)
	GetAllOffersByStatusAndOrderID(status entity.OfferStatus, orderID uint64) ([]entity.Offer, error)
}

// CreateOffer creates a new user in the database
func (a *OfferApplication) CreateOffer(balance *entity.Offer) (*entity.Offer, error) {
	return a.offerRepo.CreateOffer(balance)
}

func (a *OfferApplication) UpdateOfferByID(id uint64, offer *entity.Offer) (*entity.Offer, error) {
	return a.offerRepo.UpdateOfferByID(id, offer)
}

func (a *OfferApplication) CountOffersByDriverIDAndStatus(driverID uint64, status entity.OfferStatus) (int64, error) {
	return a.offerRepo.CountOffersByDriverIDAndStatus(driverID, status)
}

func (a *OfferApplication) CountOffersByStatusAndUserID(status entity.OfferStatus, userID uint64) (int64, error) {
	return a.offerRepo.CountOffersByStatusAndUserID(status, userID)
}

func (a *OfferApplication) GetAllOffersByStatusAndUserID(status entity.OfferStatus, userID uint64, page int, perPage int) ([]entity.Offer, error) {
	return a.offerRepo.GetAllOffersByStatusAndUserID(status, userID, page, perPage)
}

func (a *OfferApplication) GetOfferByIDAndUserID(id uint64, userID uint64) (*entity.Offer, error) {
	return a.offerRepo.GetOfferByIDAndUserID(id, userID)
}

func (a *OfferApplication) GetAllOffersByStatusAndOrderID(status entity.OfferStatus, orderID uint64) ([]entity.Offer, error) {
	return a.offerRepo.GetAllOffersByStatusAndOrderID(status, orderID)
}
