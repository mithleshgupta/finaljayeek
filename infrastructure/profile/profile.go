package profile

import (
	"context"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

// ProfileServiceInterface defines the methods that a chat service should implement.
type ProfileServiceInterface interface {
	GetProfileDetails(ctx *gin.Context, user *entity.User) (*ProfileDetails, error)
	UpdateProfileDetails(ctx *gin.Context, user *entity.User) (*ProfileDetails, error)
	UpdatePhoneNumber(ctx *gin.Context, user *entity.User) (*ProfileDetails, error)
}

// ProfileService represents the profile service implementation
type ProfileService struct {
	UserApp application.UserApplicationInterface
}

// ProfileDetails represents the details of a profile.
type ProfileDetails struct {
	LocationID uint64  `json:"location_id" validate:"required,numeric"`
	Name       string  `json:"name" validate:"required"`
	Email      *string `json:"email" validate:"omitempty,email"`
	Phone      string  `json:"phone" validate:"required,e164"`
	BirthYear  *int64  `json:"birth_year"`
}

var _ ProfileServiceInterface = &ProfileService{}

// NewProfileService creates and returns a new instance of ProfileService
func NewProfileService(userApp application.UserApplicationInterface) *ProfileService {
	return &ProfileService{
		UserApp: userApp,
	}
}

func (s *ProfileService) GetProfileDetails(ctx *gin.Context, user *entity.User) (*ProfileDetails, error) {
	return &ProfileDetails{
		LocationID: user.LocationID,
		Name:       user.Name,
		Phone:      user.Phone,
		BirthYear:  user.BirthYear,
		Email:      user.Email,
	}, nil
}

func (s *ProfileService) UpdateProfileDetails(ctx *gin.Context, user *entity.User) (*ProfileDetails, error) {
	updatedUser, err := s.UserApp.UpdateUserByID(user.ID, user)
	if err != nil {
		return nil, err
	}

	return &ProfileDetails{
		LocationID: updatedUser.LocationID,
		Name:       updatedUser.Name,
		Phone:      updatedUser.Phone,
		BirthYear:  updatedUser.BirthYear,
		Email:      updatedUser.Email,
	}, nil
}

func (s *ProfileService) UpdatePhoneNumber(ctx *gin.Context, user *entity.User) (*ProfileDetails, error) {
	updatedUser, err := s.UserApp.UpdateUserByID(user.ID, user)
	if err != nil {
		return nil, err
	}

	return &ProfileDetails{
		LocationID: updatedUser.LocationID,
		Name:       updatedUser.Name,
		Phone:      updatedUser.Phone,
		BirthYear:  updatedUser.BirthYear,
		Email:      updatedUser.Email,
	}, nil
}
