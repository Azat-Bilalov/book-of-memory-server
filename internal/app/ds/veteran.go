package ds

import "time"

type Veteran struct {
	VeteranID  uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"not null"`
	CreatedAt  time.Time
}