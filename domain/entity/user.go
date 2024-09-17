package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/config"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/security"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID                         uint64         `gorm:"primary_key;auto_increment" json:"id"`
	LocationID                 uint64         `gorm:"index;" json:"location_id" validate:"required,numeric"`
	Name                       string         `gorm:"size:255;not null;" json:"name" validate:"required"`
	Email                      *string        `gorm:"size:255;default:null" json:"email" validate:"omitempty,email"`
	Phone                      string         `gorm:"size:45;not null;unique" json:"phone" validate:"required,e164"`
	Password                   string         `gorm:"size:255;not null" json:"password" validate:"required"`
	BirthYear                  *int64         `gorm:"default:null" json:"omitempty,birth_year"`
	PhoneVerifiedAt            time.Time      `gorm:"default:null" json:"phone_verified_at"`
	Avatar                     string         `gorm:"size:255;default:null" json:"avatar"`
	Settings                   datatypes.JSON `gorm:"type:json" json:"settings"`
	CreatedAt                  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                  time.Time      `gorm:"default:null" json:"updated_at"`
	Role                       Role           `gorm:"size:255;default:user" json:"role" validate:"required,oneof=admin user"`
	Location                   Location       `gorm:"foreignKey:LocationID" json:"location"`
	IsDriver                   bool           `gorm:"-" json:"is_driver"`
	PurchasesCount             int64          `gorm:"-" json:"purchases_count"`
	OrdersCount                int64          `gorm:"-" json:"orders_count"`
	InProgressOrdersCount      int64          `gorm:"-" json:"in_progress_orders_count"`
	TripsCount                 int64          `gorm:"-" json:"trips_count"`
	BalancesSumBalance         int64          `gorm:"-" json:"balances_sum_balance"`
	MonthlyRevenue             int64          `gorm:"-" json:"monthly_revenue"`
	RatingsAvgScore            float64        `gorm:"-" json:"ratings_avg_score"`
	FastRatingsAvgScore        float64        `gorm:"-" json:"fast_ratings_avg_score"`
	ExperienceRatingsAvgScore  float64        `gorm:"-" json:"experience_ratings_avg_score"`
	RecommendedRatingsAvgScore float64        `gorm:"-" json:"recommended_ratings_avg_score"`
	ReviewsCount               int64          `gorm:"-" json:"reviews_count"`
}

type UserPublicData struct {
	ID                         uint64                   `json:"id"`
	LocationID                 uint64                   `json:"location_id"`
	Name                       string                   `json:"name"`
	Email                      *string                  `json:"email"`
	Phone                      string                   `json:"phone"`
	BirthYear                  *int64                   `json:"birth_year"`
	Avatar                     string                   `json:"avatar"`
	Settings                   []map[string]interface{} `json:"settings"`
	Location                   *LocationPublicData      `json:"location"`
	IsDriver                   bool                     `json:"is_driver"`
	PurchasesCount             int64                    `json:"purchases_count"`
	OrdersCount                int64                    `json:"orders_count"`
	InProgressOrdersCount      int64                    `json:"in_progress_orders_count"`
	TripsCount                 int64                    `json:"trips_count"`
	BalancesSumBalance         int64                    `json:"balances_sum_balance"`
	MonthlyRevenue             int64                    `json:"monthly_revenue"`
	RatingsAvgScore            float64                  `json:"ratings_avg_score"`
	FastRatingsAvgScore        float64                  `json:"fast_ratings_avg_score"`
	ExperienceRatingsAvgScore  float64                  `json:"experience_ratings_avg_score"`
	RecommendedRatingsAvgScore float64                  `json:"recommended_ratings_avg_score"`
	ReviewsCount               int64                    `json:"reviews_count"`
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// PublicData returns a copy of the user's public information
func (u *User) PublicData(languageCode string, currentUserID ...uint64) interface{} {
	conf := config.NewConfig()

	baseImageURL := conf.BaseStorageURL

	locationPublicData := u.Location.PublicData(languageCode).(*LocationPublicData)
	settingsArray := make([]map[string]interface{}, 0)
	err := json.Unmarshal(u.Settings, &settingsArray)
	if err != nil {
		// handle error
	}

	var avatar string
	if len(u.Avatar) == 0 {
		avatar = ""
	} else {
		avatar = baseImageURL + "/" + u.Avatar
	}

	return &UserPublicData{
		ID:                         u.ID,
		LocationID:                 u.LocationID,
		Name:                       u.Name,
		Email:                      u.Email,
		Phone:                      u.getPhoneNumber(currentUserID...),
		BirthYear:                  u.BirthYear,
		Avatar:                     avatar,
		Settings:                   settingsArray,
		Location:                   locationPublicData,
		IsDriver:                   u.IsDriver,
		PurchasesCount:             u.PurchasesCount,
		OrdersCount:                u.OrdersCount,
		InProgressOrdersCount:      u.InProgressOrdersCount,
		TripsCount:                 u.TripsCount,
		BalancesSumBalance:         u.BalancesSumBalance,
		MonthlyRevenue:             u.MonthlyRevenue,
		RatingsAvgScore:            u.RatingsAvgScore,
		FastRatingsAvgScore:        u.FastRatingsAvgScore,
		ExperienceRatingsAvgScore:  u.ExperienceRatingsAvgScore,
		RecommendedRatingsAvgScore: u.RecommendedRatingsAvgScore,
		ReviewsCount:               u.ReviewsCount,
	}
}

func (u *User) getPhoneNumber(currentUserID ...uint64) string {
	if len(currentUserID) > 0 && currentUserID[0] == u.ID {
		return u.Phone
	}
	return ""
}

// BeforeSave is a gorm hook that hashes the user's password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Check if the password is already hashed
	if !security.IsHashed(u.Password) {
		print("Password is hashed")
		print(u.Password)
		// Hash the password using the security.Hash function
		hash, err := security.Hash(u.Password)
		if err != nil {
			return err
		}

		// Set the "password" column to the hashed password
		tx.Statement.SetColumn("password", string(hash))
	}

	return nil
}

// AfterFind is a gorm hook that sets the value of the IsDriver field
// based on whether a driver with the same user_id exists
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	var driversCount int64
	err = tx.Table("drivers").Where("user_id = ?", u.ID).Count(&driversCount).Error
	if err != nil {
		return err
	}

	if driversCount > 0 {
		// Set the value of the IsDriver field
		u.IsDriver = true
	} else {
		// Set the value of the IsDriver field
		u.IsDriver = false
	}

	var purchasesCount int64
	err = tx.Table("orders").Where("user_id = ?", u.ID).Count(&purchasesCount).Error
	if err != nil {
		return err
	}
	u.PurchasesCount = purchasesCount

	// Driver struct to store the retrieved driver data
	var driver Driver

	// Retrieve the driver data using the user ID
	if err := tx.Table("drivers").Where("user_id = ?", u.ID).First(&driver).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle the error appropriately, excluding the "record not found" error
			return err
		}

		// Handle the "record not found" case separately
		driver = Driver{}
	}

	var ordersCount int64

	if driver.ID != 0 {
		// Count the number of orders associated with the driver
		if err := tx.Table("orders").Where("driver_id = ?", driver.ID).Count(&ordersCount).Error; err != nil {
			// Handle the error appropriately
			return err
		}
		u.OrdersCount = ordersCount
	} else {
		u.OrdersCount = 0
	}

	var inProgressOffersCount, inProgressOrdersCount int64

	if driver.ID != 0 {
		if err := tx.Table("offers").Where("driver_id = ?", driver.ID).Where("status IN (?)", OfferStatusPending).Count(&inProgressOffersCount).Error; err != nil {
			// Handle the error appropriately
			return err
		}
		if err := tx.Table("orders").Where("driver_id = ?", driver.ID).Where("status NOT IN (?, ?, ?)", OrderCompletedStatus, OrderCanceledStatus, ShipmentReturnedStatus).Count(&inProgressOrdersCount).Error; err != nil {
			// Handle the error appropriately
			return err
		}
		u.InProgressOrdersCount = inProgressOffersCount + inProgressOrdersCount
	} else {
		u.InProgressOrdersCount = 0
	}

	var tripsCount int64

	if driver.ID != 0 {
		// Count the number of orders associated with the driver
		if err := tx.Table("orders").Where("driver_id = ?", driver.ID).Where("status IN (?, ?)", ShipmentDeliveredStatus, OrderCompletedStatus).Count(&tripsCount).Error; err != nil {
			// Handle the error appropriately
			return err
		}
		u.TripsCount = tripsCount
	} else {
		u.TripsCount = 0
	}

	if driver.ID != 0 {
		var balancesSumBalance int64
		var monthlyRevenue int64

		// Calculate the sum of balances associated with the driver
		result := tx.Table("balances").Select("COALESCE(SUM(balance), 0)").Where("driver_id = ?", driver.ID).Row()
		if err := result.Scan(&balancesSumBalance); err != nil {
			// Handle the error appropriately
			return err
		}
		u.BalancesSumBalance = balancesSumBalance

		result = tx.Table("balances").Select("COALESCE(SUM(balance), 0)").
			Where("driver_id = ? AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_DATE)", driver.ID).
			Row()
		if err := result.Scan(&monthlyRevenue); err != nil {
			// Handle the error appropriately
			return err
		}
		u.MonthlyRevenue = monthlyRevenue
	} else {
		u.BalancesSumBalance = 0
		u.MonthlyRevenue = 0
	}

	if driver.ID != 0 {
		var ratingsAvgScore float64
		var fastRatingsAvgScore float64
		var experienceRatingsAvgScore float64
		var recommendedRatingsAvgScore float64

		subQuery := fmt.Sprintf("SELECT id FROM orders WHERE driver_id = %d", driver.ID)

		err := tx.Raw(fmt.Sprintf(`
			SELECT LEAST(COALESCE((AVG(fast_rating) + AVG(experience_rating) + AVG(recommended_rating)) / 3, 0), 5)
			FROM ratings
			WHERE order_id IN (%s)
		`, subQuery)).Scan(&ratingsAvgScore).Error
		if err != nil {
			return err
		}

		err = tx.Raw(fmt.Sprintf(`
			SELECT LEAST(COALESCE(AVG(fast_rating), 0), 5)
			FROM ratings
			WHERE order_id IN (%s)
		`, subQuery)).Scan(&fastRatingsAvgScore).Error
		if err != nil {
			return err
		}

		err = tx.Raw(fmt.Sprintf(`
			SELECT LEAST(COALESCE(AVG(experience_rating), 0), 5)
			FROM ratings
			WHERE order_id IN (%s)
		`, subQuery)).Scan(&experienceRatingsAvgScore).Error
		if err != nil {
			return err
		}

		err = tx.Raw(fmt.Sprintf(`
			SELECT LEAST(COALESCE(AVG(recommended_rating), 0), 5)
			FROM ratings
			WHERE order_id IN (%s)
		`, subQuery)).Scan(&recommendedRatingsAvgScore).Error
		if err != nil {
			return err
		}

		// Format the floating-point number with one decimal place
		formattedRating := fmt.Sprintf("%.1f", ratingsAvgScore)

		// Convert the formatted string back to float64
		formattedRatingFloat, err := strconv.ParseFloat(formattedRating, 64)
		if err != nil {
			return err
		}

		// Format the floating-point number with one decimal place
		formattedFastRating := fmt.Sprintf("%.1f", fastRatingsAvgScore)

		// Convert the formatted string back to float64
		formattedFastRatingFloat, err := strconv.ParseFloat(formattedFastRating, 64)
		if err != nil {
			return err
		}

		// Format the floating-point number with one decimal place
		formattedExperienceRating := fmt.Sprintf("%.1f", experienceRatingsAvgScore)

		// Convert the formatted string back to float64
		formattedExperienceRatingFloat, err := strconv.ParseFloat(formattedExperienceRating, 64)
		if err != nil {
			return err
		}

		// Format the floating-point number with one decimal place
		formattedRecommendedRating := fmt.Sprintf("%.1f", recommendedRatingsAvgScore)

		// Convert the formatted string back to float64
		formattedRecommendedRatingFloat, err := strconv.ParseFloat(formattedRecommendedRating, 64)
		if err != nil {
			return err
		}

		u.RatingsAvgScore = formattedRatingFloat
		u.FastRatingsAvgScore = formattedFastRatingFloat
		u.ExperienceRatingsAvgScore = formattedExperienceRatingFloat
		u.RecommendedRatingsAvgScore = formattedRecommendedRatingFloat
	} else {
		u.RatingsAvgScore = 0
		u.FastRatingsAvgScore = 0
		u.ExperienceRatingsAvgScore = 0
		u.RecommendedRatingsAvgScore = 0
	}

	var reviewsCount int64

	if driver.ID != 0 {
		subQuery := fmt.Sprintf("SELECT id FROM orders WHERE driver_id = %d", driver.ID)

		err := tx.Raw(fmt.Sprintf(`
		SELECT COALESCE(COUNT(ratings.id), 0)
		FROM ratings
		WHERE order_id IN (%s)
		`, subQuery)).Scan(&reviewsCount).Error
		if err != nil {
			return err
		}

		u.ReviewsCount = reviewsCount
	} else {
		u.ReviewsCount = 0
	}

	return
}

func (u *User) AddSetting(key string, value interface{}) error {
	if u.Settings == nil {
		u.Settings = datatypes.JSON([]byte("{}"))
	}

	settingsMap := make(map[string]interface{})
	err := json.Unmarshal(u.Settings, &settingsMap)
	if err != nil {
		return err
	}

	settingsMap[key] = value

	updatedSettings, err := json.Marshal(settingsMap)
	if err != nil {
		return err
	}

	u.Settings = datatypes.JSON(updatedSettings)
	return nil
}

// GetSettingByKey retrieves the value of a specific setting based on the given key
func (u *User) GetSettingByKey(key string) (interface{}, error) {
	if u.Settings == nil {
		return nil, errors.New("settings field is empty")
	}

	settingsMap := make(map[string]interface{})
	err := json.Unmarshal(u.Settings, &settingsMap)
	if err != nil {
		return nil, err
	}

	value, ok := settingsMap[key]
	if !ok {
		return nil, errors.New("setting not found")
	}

	return value, nil
}
