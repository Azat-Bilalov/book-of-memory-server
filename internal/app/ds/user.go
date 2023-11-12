package ds

import (
	"time"
)

const (
	USER_ROLE_MODERATOR = "moderator"
	USER_ROLE_USER      = "user"
)

type User struct {
	User_id   string `gorm:"primaryKey"`
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

type UserResponse struct {
	User_id   string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
