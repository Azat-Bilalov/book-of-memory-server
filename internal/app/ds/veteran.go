package ds

import (
	"mime/multipart"
	"time"
)

type Veteran struct {
	Veteran_id string    `gorm:"primarykey;default:gen_random_uuid()" json:"veteran_id"`
	FirstName  string    `gorm:"not null" json:"first_name"`
	LastName   string    `gorm:"not null" json:"last_name"`
	Patronymic string    `gorm:"not null" json:"patronymic"`
	BirthDate  time.Time `gorm:"not null" json:"birth_date"`
	ImageUrl   string    `gorm:"not null" json:"image_url"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	// Documents  []Document `gorm:"foreignKey:Veteran_id;references:Veteran_id" json:"documents"`
}

type VeteranRequest struct {
	FirstName  string                `form:"first_name" veteran:"required"`
	LastName   string                `form:"last_name" veteran:"required"`
	Patronymic string                `form:"patronymic" veteran:"required"`
	BirthDate  time.Time             `form:"birth_date" veteran:"required" time_format:"2006-01-02" time_utc:"true" time_location:"Europe/Moscow"`
	Image      *multipart.FileHeader `form:"image" veteran:"required" swaggerignore:"true"`
}
