package ds

import "time"

type Binding struct {
	Binding_id  uint   `gorm:"primaryKey"`
	Status      string `gorm:"not null;check:status IN ('entered', 'in_progress', 'completed', 'canceled', 'deleted')"`
	Info        string `gorm:"type:text"`
	FileURL     string
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	FormattedAt time.Time
	EndedAt     time.Time
	UserID      uint `gorm:"not null"`
	ModeratorID uint `gorm:"not null"`
	VeteranID   uint `gorm:"not null"`
}
