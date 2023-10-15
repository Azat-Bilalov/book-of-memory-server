package ds

import "time"

type User struct {
	User_id   uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Passwd    string `gorm:"not null"`
	Role      string `gorm:"not null;default:user;check:role IN ('moderator', 'user')"`
	CreatedAt time.Time
}

type UserCreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}
