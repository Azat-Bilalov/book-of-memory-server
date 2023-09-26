package ds

import "time"

type User struct {
	UserID    uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Passwd    string `gorm:"not null"`
	Role      string `gorm:"not null;default:user;check:role IN ('moderator', 'user')"`
	CreatedAt time.Time
}
