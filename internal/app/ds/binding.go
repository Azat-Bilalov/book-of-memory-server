package ds

import "time"

const (
	BINDING_STATUS_ENTERED     = "entered"
	BINDING_STATUS_IN_PROGRESS = "in_progress"
	BINDING_STATUS_COMPLETED   = "completed"
	BINDING_STATUS_CANCELED    = "canceled"
	BINDING_STATUS_DELETED     = "deleted"
)

type Binding struct {
	Binding_id  string    `gorm:"primarykey;default:gen_random_uuid()"`
	Status      string    `gorm:"not null;check:status IN ('entered', 'in_progress', 'completed', 'canceled', 'deleted');default:'entered"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	FormattedAt *time.Time
	EndedAt     *time.Time
	UserID      string     `gorm:"not null"`
	ModeratorID string     `gorm:"not null"`
	VeteranID   *string    `gorm:"not null"`
	Veteran     *Veteran   `gorm:"foreignKey:veteran_id;references:veteran_id"`
	Documents   []Document `gorm:"many2many:doc_bindings;foreignKey:binding_id;joinForeignKey:binding_id;References:document_id;JoinReferences:document_id"`
}

type BindingUpdateRequest struct {
	VeteranID string `json:"veteran_id"`
}
