package ds

import (
	"time"
)

const (
	USER_ROLE_MODERATOR = "moderator"
	USER_ROLE_USER      = "user"
)

type User struct {
	User_id   string `gorm:"primaryKey" json:"user_id"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"unique" json:"email"`
	Passwd    string `gorm:"not null" json:"passwd"`
	Role      string `gorm:"not null;default:user;check:role IN ('moderator', 'user')" json:"role"`
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
