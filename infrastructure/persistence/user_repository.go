package persistence

import (
	"fmt"

	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/security"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository implements repository.UserRepository
// and handles CRUD operations for User entities
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if err := r.db.Debug().Model(&user).Create(&user).Error; err != nil {
		return nil, err
	}

	if err := r.db.Debug().Model(&user).Preload("Location").Take(&user).Error; err != nil {
		return nil, err
	}

	var orders []entity.Order
	if err := r.db.Debug().Table("orders").
		Where("recipient_phone_number = ?", user.Phone).
		Where("status != ?", entity.OrderCanceledStatus).
		Where("status != ?", entity.OrderCompletedStatus).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	for _, order := range orders {
		if err := r.db.Debug().Model(&order).Updates(entity.Order{
			RecipientID: user.ID,
		}).Where("id = ?", order.ID).Error; err != nil {
			return nil, err
		}
	}

	return user, nil
}

// UpdateUserByID updates the user
func (r *UserRepository) UpdateUserByID(id uint64, user *entity.User) (*entity.User, error) {
	if err := r.db.Debug().Model(&user).Updates(user).Where("id = ?", id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UserWithFieldExists checks if a user with the given field and value exists in the database
func (r *UserRepository) UserWithFieldExists(field string, value string) (bool, error) {
	var user entity.User
	var usersCount int64
	if err := r.db.Debug().Model(&user).Select("id").Where(fmt.Sprintf("%s = ?", field), value).Count(&usersCount).Error; err != nil {
		return false, err
	}
	return usersCount > 0, nil
}

// GetUserByPhoneAndPassword retrieves a user from the database by phone and password
func (r *UserRepository) GetUserByPhoneAndPassword(phone string, password string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Debug().Model(&user).Where("phone = ?", phone).Preload("Location").Take(&user).Error; err != nil {
		return nil, err
	}
	// Verify the password
	if err := security.VerifyPassword(user.Password, password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("incorrect password")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	if err := r.db.Debug().Model(&entity.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetAllUsers(page int, perPage int) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.Debug().Model(&entity.User{}).Limit(perPage).Offset((page - 1) * perPage).Preload("Location").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id uint64) (*entity.User, error) {
	var user entity.User
	if err := r.db.Debug().Model(&user).Where("id = ?", id).Preload("Location").Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByPhone(phone string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Debug().Model(&user).Where("phone = ?", phone).Preload("Location").Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
