package persistence

import (
	"fmt"
	"time"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"gorm.io/gorm"
)

// OrderRepository implements the repository.OrderRepository interface
type OrderRepository struct {
	// db is a pointer to the GORM DB instance
	db *gorm.DB
}

// NewOrderRepository creates a new instance of the OrderRepository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order in the database
func (r *OrderRepository) CreateOrder(order *entity.Order) (*entity.Order, error) {
	if err := r.db.Debug().Model(&order).Create(&order).Error; err != nil {
		return nil, err
	}

	if err := r.db.Debug().Model(&order).Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Take(&order).Error; err != nil {
		return nil, err
	}

	var drivers []entity.Driver
	if err := r.db.Debug().Table("drivers").Select("drivers.*").
		Where("ST_DWithin(ST_MakePoint(drivers.longitude, drivers.latitude)::geography, ST_MakePoint(?, ?)::geography, ?)", order.Longitude, order.Latitude, 50000).
		Order(fmt.Sprintf("ST_Distance(ST_MakePoint(drivers.longitude, drivers.latitude)::geography, ST_MakePoint(%.6f, %.6f)::geography)", order.Longitude, order.Latitude)).
		Preload("User").Preload("User.Location").Find(&drivers).Error; err != nil {
		return nil, err
	}

	for _, driver := range drivers {
		if isAvailableSetting, err := driver.User.GetSettingByKey("is_available"); err == nil {
			isAvailable, ok := isAvailableSetting.(bool)
			if !ok {
				isAvailable = false // Default value if the setting is not a boolean
			}

			if isAvailable {
				var orderDriverPool entity.OrderDriverPool
				orderDriverPool.DriverID = driver.ID
				orderDriverPool.OrderID = order.ID

				if err := r.db.Debug().Model(&orderDriverPool).Create(&orderDriverPool).Error; err != nil {
					// Handle the error internally (e.g., log the error)
					print(err)
				}
			}
		} else {
			// Handle the error internally (e.g., log the error)
		}
	}

	return order, nil
}

// UpdateOrder updates the order
func (r *OrderRepository) UpdateOrderByID(id uint64, order *entity.Order) (*entity.Order, error) {
	if err := r.db.Debug().Model(&order).Updates(order).Where("id = ?", id).Error; err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderDriverPool updates the order driver pool
func (r *OrderRepository) UpdateOrderDriverPoolByOrderIDAndDriverID(orderID uint64, driverID uint64, orderDriverPool *entity.OrderDriverPool) (*entity.OrderDriverPool, error) {
	if err := r.db.Debug().Model(&orderDriverPool).Where("order_id = ?", orderID).Where("driver_id = ?", driverID).Updates(orderDriverPool).Error; err != nil {
		return nil, err
	}

	return orderDriverPool, nil
}

func (r *OrderRepository) CountOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus) (int64, error) {
	var count int64
	if err := r.db.Debug().Where("user_id = ?", userID).Where("status NOT IN (?)", status).Model(&entity.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) CountOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus) (int64, error) {
	var count int64
	if err := r.db.Debug().Where("(user_id = ? OR recipient_id = ?) AND status NOT IN (?)", userID, recipientID, status).Model(&entity.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) CountOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus) (int64, error) {
	var count int64
	if err := r.db.Debug().Where("driver_id = ?", driverID).Where("status NOT IN (?)", status).Model(&entity.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetAllOrdersByUserIDExcludingStatus(userID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Debug().Where("user_id = ?", userID).Where("status NOT IN (?)", status).Model(&entity.Order{}).Limit(perPage).Offset((page - 1) * perPage).Order("created_at desc").Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetAllOrdersByUserIDAndRecipientIDExcludingStatus(userID uint64, recipientID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Debug().Where("(user_id = ? OR recipient_id = ?) AND status NOT IN (?)", userID, recipientID, status).Model(&entity.Order{}).Limit(perPage).Offset((page - 1) * perPage).Order("created_at desc").Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetAllOrdersByDriverIDExcludingStatus(driverID uint64, status []entity.OrderStatus, page int, perPage int) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Debug().Where("driver_id = ?", driverID).Where("status NOT IN (?)", status).Model(&entity.Order{}).Limit(perPage).Offset((page - 1) * perPage).Order("created_at desc").Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetAllOrdersByRecipientPhoneNumber(recipientPhoneNumber string) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Debug().Where("recipient_phone_number = ?", recipientPhoneNumber).Model(&entity.Order{}).Order("created_at desc").Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) CountDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, byArrival *string) (int64, error) {
	var count int64
	db := r.db.Debug().Table("orders").
		Joins("JOIN order_driver_pools ON orders.id = order_driver_pools.order_id").
		Joins("JOIN drivers ON order_driver_pools.driver_id = drivers.id").
		Joins("JOIN delivery_times ON orders.delivery_time_id = delivery_times.id").
		Joins("JOIN locations AS source ON orders.location_id = source.id").
		Joins("JOIN locations AS destination ON orders.destination_id = destination.id").
		Where("order_driver_pools.driver_id = ?", driverID).
		Where("order_driver_pools.status = ?", entity.PendingStatus).
		Where("orders.category_id = ?", categoryID)

	// Handle arrival filtering if byArrival is provided
	if byArrival != nil && (*byArrival == "today" || *byArrival == "tomorrow") {
		today := time.Now().Truncate(24 * time.Hour)

		// Filter the results based on the arrival time
		if *byArrival == "today" {
			db = db.Where("(current_timestamp + (delivery_times.duration * interval '1 second')) >= ?", today.Add(24*time.Hour)).Where("(current_timestamp + (delivery_times.duration * interval '1 second')) < ?", today.Add(48*time.Hour))
		} else if *byArrival == "tomorrow" {
			db = db.Where("(current_timestamp + (delivery_times.duration * interval '1 second')) >= ?", today).Where("(current_timestamp + (delivery_times.duration * interval '1 second')) < ?", today.Add(24*time.Hour))
		}
	}

	err := db.Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) GetAllDriverPoolsByDriverIDAndCategoryID(driverID uint64, categoryID uint64, page int, perPage int, orderBy *string, byArrival *string) ([]entity.Order, error) {
	// Build the base query
	db := r.db.Debug().Table("orders").
		Select("orders.*").
		Joins("JOIN order_driver_pools ON orders.id = order_driver_pools.order_id").
		Joins("JOIN drivers ON order_driver_pools.driver_id = drivers.id").
		Joins("JOIN delivery_times ON orders.delivery_time_id = delivery_times.id").
		Joins("JOIN locations AS source ON orders.location_id = source.id").
		Joins("JOIN locations AS destination ON orders.destination_id = destination.id").
		Where("order_driver_pools.driver_id = ?", driverID).
		Where("order_driver_pools.status = ?", entity.PendingStatus).
		Where("orders.category_id = ?", categoryID)

	// Handle sorting if orderBy is provided
	if orderBy != nil && (*orderBy == "distance" || *orderBy == "arrival" || *orderBy == "created_at") {
		switch *orderBy {
		case "distance":
			// Include sorting logic based on distance
			db = db.Order("SQRT(POW(orders.latitude - destination.latitude, 2) + POW(orders.longitude - destination.longitude, 2)) ASC")
		case "arrival":
			// Include sorting logic based on expected arrival time
			db = db.Order("(current_timestamp + (delivery_times.duration * interval '1 second')) ASC")
		case "created_at":
			// Include sorting logic based on creation time
			db = db.Order("orders.created_at DESC")
		}
	}

	// Handle arrival filtering if byArrival is provided
	if byArrival != nil && (*byArrival == "today" || *byArrival == "tomorrow") {
		today := time.Now().Truncate(24 * time.Hour)

		// Filter the results based on the arrival time
		if *byArrival == "today" {
			db = db.Where("(current_timestamp + (delivery_times.duration * interval '1 second')) >= ?", today).Where("(current_timestamp + (delivery_times.duration * interval '1 second')) < ?", today.Add(24*time.Hour))
		} else if *byArrival == "tomorrow" {
			db = db.Where("(current_timestamp + (delivery_times.duration * interval '1 second')) >= ?", today.Add(24*time.Hour)).Where("(current_timestamp + (delivery_times.duration * interval '1 second')) < ?", today.Add(48*time.Hour))
		}
	}

	// Paginate the results
	offset := (page - 1) * perPage
	db = db.Offset(offset).Limit(perPage).
		Preload("Location").
		Preload("User").
		Preload("User.Location").
		Preload("Driver").
		Preload("Driver.User").
		Preload("Driver.User.Location").
		Preload("Driver.TransportationMode").
		Preload("Recipient").
		Preload("Recipient.Location").
		Preload("Category").
		Preload("Size").
		Preload("DeliveryTime").
		Preload("ShipmentContents").
		Preload("ExtraServices").
		Preload("Destination")

	var orders []entity.Order
	err := db.Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID retrieves a order by its ID
func (r *OrderRepository) GetOrderByID(id uint64) (*entity.Order, error) {
	// Order struct to store the retrieved order data
	var order entity.Order
	// Find the order by its ID and store the data in the order struct
	if err := r.db.Debug().Where("id = ?", id).Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Take(&order).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the order data and nil error
	return &order, nil
}

func (r *OrderRepository) GetOrderByIDAndUserID(id uint64, userID uint64) (*entity.Order, error) {
	// Order struct to store the retrieved order data
	var order entity.Order
	// Find the order by its ID and store the data in the order struct
	if err := r.db.Debug().Where("id = ?", id).Where("user_id = ?", userID).Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Take(&order).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the order data and nil error
	return &order, nil
}

func (r *OrderRepository) GetOrderByIDAndRecipientID(id uint64, recipientID uint64) (*entity.Order, error) {
	// Order struct to store the retrieved order data
	var order entity.Order
	// Find the order by its ID and store the data in the order struct
	if err := r.db.Debug().Where("id = ?", id).Where("recipient_id = ?", recipientID).Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Take(&order).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the order data and nil error
	return &order, nil
}

func (r *OrderRepository) GetOrderByIDAndDriverID(id uint64, driverID uint64) (*entity.Order, error) {
	// Order struct to store the retrieved order data
	var order entity.Order
	// Find the order by its ID and store the data in the order struct
	if err := r.db.Debug().Where("id = ?", id).Where("driver_id = ?", driverID).Preload("Location").Preload("User").Preload("User.Location").Preload("Driver").Preload("Driver.User").Preload("Driver.User.Location").Preload("Driver.TransportationMode").Preload("Recipient").Preload("Recipient.Location").Preload("Category").Preload("Size").Preload("DeliveryTime").Preload("ShipmentContents").Preload("ExtraServices").Preload("Destination").Take(&order).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the order data and nil error
	return &order, nil
}

func (r *OrderRepository) GetOrderDriverPoolByOrderIDAndDriverID(orderID uint64, driverID uint64) (*entity.OrderDriverPool, error) {
	// Order driver pool struct to store the retrieved order driver pool data
	var orderDriverPool entity.OrderDriverPool
	if err := r.db.Debug().Where("order_id = ?", orderID).Where("driver_id = ?", driverID).Take(&orderDriverPool).Error; err != nil {
		// If there's an error, return nil and the error
		return nil, err
	}
	// return the order driver pool data and nil error
	return &orderDriverPool, nil
}
