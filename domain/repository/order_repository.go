package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// OrderRepository defines the methods for interacting with order data
type OrderRepository interface {
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
	CountDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, byArrival *string) (int64, error)
	GetAllDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, page int, perPage int, orderBy *string, byArrival *string) ([]entity.Order, error)
	GetOrderByID(uint64) (*entity.Order, error)
	GetOrderByIDAndUserID(id uint64, userID uint64) (*entity.Order, error)
	GetOrderByIDAndRecipientID(id uint64, recipientID uint64) (*entity.Order, error)
	GetOrderByIDAndDriverID(id uint64, driverID uint64) (*entity.Order, error)
	GetOrderDriverPoolByOrderIDAndDriverID(uint64, uint64) (*entity.OrderDriverPool, error)
}
