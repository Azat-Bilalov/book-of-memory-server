package ds

import "time"

type Veteran struct {
	Veteran_id string `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"not null"`
	ImageUrl   string
	CreatedAt  time.Time
}
