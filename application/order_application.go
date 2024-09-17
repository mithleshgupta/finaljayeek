package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// OrderApplication handles the business logic for orders
type OrderApplication struct {
	orderRepo repository.OrderRepository
}

var _ OrderApplicationInterface = &OrderApplication{}

// OrderApplicationInterface defines the methods available for OrderApplication
type OrderApplicationInterface interface {
	CreateOrder(order *entity.Order) (*entity.Order, error)
	UpdateOrderByID(id uint64, order *entity.Order) (*entity.Order, error)
	UpdateOrderDriverPoolByOrderIDAndDriverID(orderID uint64, driverID uint64, orderDriverPool *entity.OrderDriverPool) (*entity.OrderDriverPool, error)
	CountOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus) (int64, error)
	CountOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus) (int64, error)
	CountOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus) (int64, error)
	GetAllOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error)
	GetAllOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error)
	GetAllOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error)
	GetAllOrdersByRecipientPhoneNumber(recipientPhoneNumber string) ([]entity.Order, error)
	CountDriverPoolsByDriverIDAndCategoryID(uint64, uint64, *string) (int64, error)
	GetAllDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, page int, perPage int, orderBy *string, byArrival *string) ([]entity.Order, error)
	GetOrderByID(uint64) (*entity.Order, error)
	GetOrderByIDAndUserID(id uint64, userID uint64) (*entity.Order, error)
	GetOrderByIDAndRecipientID(id uint64, recipientID uint64) (*entity.Order, error)
	GetOrderByIDAndDriverID(id uint64, driverID uint64) (*entity.Order, error)
	GetOrderDriverPoolByOrderIDAndDriverID(uint64, uint64) (*entity.OrderDriverPool, error)
}

// CreateOrder creates a new order in the database
func (a *OrderApplication) CreateOrder(order *entity.Order) (*entity.Order, error) {
	return a.orderRepo.CreateOrder(order)
}

func (a *OrderApplication) UpdateOrderByID(id uint64, order *entity.Order) (*entity.Order, error) {
	return a.orderRepo.UpdateOrderByID(id, order)
}

func (a *OrderApplication) UpdateOrderDriverPoolByOrderIDAndDriverID(orderID uint64, driverID uint64, orderDriverPool *entity.OrderDriverPool) (*entity.OrderDriverPool, error) {
	return a.orderRepo.UpdateOrderDriverPoolByOrderIDAndDriverID(orderID, driverID, orderDriverPool)
}

func (a *OrderApplication) CountOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus) (int64, error) {
	return a.orderRepo.CountOrdersByUserIDExcludingStatus(userID, status)
}

func (a *OrderApplication) CountOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus) (int64, error) {
	return a.orderRepo.CountOrdersByUserIDAndRecipientIDExcludingStatus(userID, recipientID, status)
}

func (a *OrderApplication) CountOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus) (int64, error) {
	return a.orderRepo.CountOrdersByDriverIDExcludingStatus(driverID, status)
}

func (a *OrderApplication) GetAllOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	return a.orderRepo.GetAllOrdersByUserIDExcludingStatus(userID, status, page, perPage)
}

func (a *OrderApplication) GetAllOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	return a.orderRepo.GetAllOrdersByUserIDAndRecipientIDExcludingStatus(userID, recipientID, status, page, perPage)
}

func (a *OrderApplication) GetAllOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	return a.orderRepo.GetAllOrdersByDriverIDExcludingStatus(driverID, status, page, perPage)
}

func (a *OrderApplication) GetAllOrdersByRecipientPhoneNumber(recipientPhoneNumber string) ([]entity.Order, error) {
	return a.orderRepo.GetAllOrdersByRecipientPhoneNumber(recipientPhoneNumber)
}

func (a *OrderApplication) CountDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, byArrival *string) (int64, error) {
	return a.orderRepo.CountDriverPoolsByDriverIDAndCategoryID(driverID, categoryID, byArrival)
}

func (a *OrderApplication) GetAllDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, page int, perPage int, orderBy *string, byArrival *string) ([]entity.Order, error) {
	return a.orderRepo.GetAllDriverPoolsByDriverIDAndCategoryID(driverID, categoryID, page, perPage, orderBy, byArrival)
}

// GetByID returns a order by its ID
func (a *OrderApplication) GetOrderByID(id uint64) (*entity.Order, error) {
	return a.orderRepo.GetOrderByID(id)
}

func (a *OrderApplication) GetOrderByIDAndUserID(id uint64, userID uint64) (*entity.Order, error) {
	return a.orderRepo.GetOrderByIDAndUserID(id, userID)
}

func (a *OrderApplication) GetOrderByIDAndRecipientID(id uint64, recipientID uint64) (*entity.Order, error) {
	return a.orderRepo.GetOrderByIDAndRecipientID(id, recipientID)
}

func (a *OrderApplication) GetOrderByIDAndDriverID(id uint64, driverID uint64) (*entity.Order, error) {
	return a.orderRepo.GetOrderByIDAndDriverID(id, driverID)
}

func (a *OrderApplication) GetOrderDriverPoolByOrderIDAndDriverID(orderID uint64, driverID uint64) (*entity.OrderDriverPool, error) {
	return a.orderRepo.GetOrderDriverPoolByOrderIDAndDriverID(orderID, driverID)
}
