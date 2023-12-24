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
	FirstName  string                `json:"first_name"`
	LastName   string                `json:"last_name"`
	Patronymic string                `json:"patronymic"`
	BirthDate  time.Time             `json:"birth_date"`
	Image      *multipart.FileHeader `form:"image" veteran:"required" swaggerignore:"true"`
}
