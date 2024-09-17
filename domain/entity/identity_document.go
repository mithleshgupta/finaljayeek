package entity

type IdentityDocument struct {
	ID                       uint64 `gorm:"primaryKey" json:"id"`
	UserID                   uint64 `json:"user_id" validate:"required,numeric"`
	DrivingLicenseIDImage    string `gorm:"not null" json:"driving_license_id_image"`
	VehicleRegistrationImage string `gorm:"not null" json:"vehicle_registration_image"`
	VehicleFrontPhotoImage   string `gorm:"not null" json:"vehicle_front_photo_image"`
	LivePhotoWithIDImage     string `gorm:"not null" json:"live_photo_with_id_image"`
	User                     User   `gorm:"foreignKey:UserID" json:"user"`
}
