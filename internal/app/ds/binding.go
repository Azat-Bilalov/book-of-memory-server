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
	Binding_id  string     `gorm:"primarykey;default:gen_random_uuid()" json:"binding_id"`
	Status      string     `gorm:"not null;check:status IN ('entered', 'in_progress', 'completed', 'canceled', 'deleted');default:'entered" json:"status"`
	CreatedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	FormattedAt *time.Time `gorm:"not null" json:"formatted_at"`
	EndedAt     *time.Time `gorm:"not null" json:"ended_at"`
	UserID      string     `gorm:"not null" json:"user_id"`
	ModeratorID string     `gorm:"not null" json:"moderator_id"`
	VeteranID   *string    `gorm:"not null" json:"veteran_id"`
	User        *User      `gorm:"hasOne:users;foreignKey:user_id;references:user_id;" json:"user"`
	Moderator   *User      `gorm:"hasOne:users;foreignKey:moderator_id;references:user_id;" json:"moderator"`
	Veteran     *Veteran   `gorm:"foreignKey:veteran_id;references:veteran_id" json:"veteran"`
	Documents   []Document `gorm:"many2many:doc_bindings;foreignKey:binding_id;joinForeignKey:binding_id;References:document_id;JoinReferences:document_id" json:"documents"`
}

type BindingUpdateRequest struct {
	VeteranID string `json:"veteran_id"`
}

type BindingStatusUpdateRequest struct {
	Status string `json:"status"`
}
