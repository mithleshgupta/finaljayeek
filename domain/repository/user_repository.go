package repository

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
)

// UserRepository defines the methods that a user repository should implement
type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUserByID(id uint64, user *entity.User) (*entity.User, error)
	UserWithFieldExists(field string, value string) (bool, error)
	GetUserByPhoneAndPassword(phone string, password string) (*entity.User, error)
	CountUsers() (int64, error)
	GetAllUsers(page int, perPage int) ([]entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
}
