package ds

import (
	"mime/multipart"
	"time"
)

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

type DocumentRequest struct {
	Title       string                `form:"title" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
}

type DocumentResponse struct {
	Document_id string `json:"document_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type DocumentListResponse struct {
	Documents         []*Document `json:"documents"`
	EnteredBinding_id *string     `json:"entered_binding_id"`
}
