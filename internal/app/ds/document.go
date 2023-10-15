package ds

import "time"

const (
	DOCUMENT_STATUS_ACTIVE  = "active"
	DOCUMENT_STATUS_DELETED = "deleted"
)

type Document struct {
	Document_id string `gorm:"primarykey;default:gen_random_uuid()"`
	Title       string
	Description string
	Image_url   string
	Status      string `gorm:"check:status IN ('active', 'deleted')"`
	CreatedAt   time.Time
}

type DocumentCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image_url   string `json:"image_url"`
}

type DocumentUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image_url   string `json:"image_url"`
}
