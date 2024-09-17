package application

import (
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/domain/repository"
)

// UserApplication struct is responsible for handling user-related business logic
type UserApplication struct {
	userRepo repository.UserRepository
}

var _ UserApplicationInterface = &UserApplication{}

// UserApplicationInterface defines the methods that UserApplication should implement
type UserApplicationInterface interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUserByID(id uint64, user *entity.User) (*entity.User, error)
	UserWithFieldExists(string, string) (bool, error)
	GetUserByPhoneAndPassword(string, string) (*entity.User, error)
	CountUsers() (int64, error)
	GetAllUsers(page int, perPage int) ([]entity.User, error)
	GetUserByID(uint64) (*entity.User, error)
	GetUserByPhone(string) (*entity.User, error)
}

// CreateUser creates a new user in the database
func (a *UserApplication) CreateUser(user *entity.User) (*entity.User, error) {
	return a.userRepo.CreateUser(user)
}

func (a *UserApplication) UpdateUserByID(id uint64, user *entity.User) (*entity.User, error) {
	return a.userRepo.UpdateUserByID(id, user)
}

// UserWithFieldExists checks if a user with the given field and value exists in the database
func (a *UserApplication) UserWithFieldExists(field string, value string) (bool, error) {
	return a.userRepo.UserWithFieldExists(field, value)
}

// GetUserByPhoneAndPassword retrieves a user from the database by phone and password
func (a *UserApplication) GetUserByPhoneAndPassword(phone string, password string) (*entity.User, error) {
	return a.userRepo.GetUserByPhoneAndPassword(phone, password)
}

func (a *UserApplication) CountUsers() (int64, error) {
	return a.userRepo.CountUsers()
}

func (a *UserApplication) GetAllUsers(page int, perPage int) ([]entity.User, error) {
	return a.userRepo.GetAllUsers(page, perPage)
}

func (a *UserApplication) GetUserByID(id uint64) (*entity.User, error) {
	return a.userRepo.GetUserByID(id)
}

func (a *UserApplication) GetUserByPhone(phone string) (*entity.User, error) {
	return a.userRepo.GetUserByPhone(phone)
}
