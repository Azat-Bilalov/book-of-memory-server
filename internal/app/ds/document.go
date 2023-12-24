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
	Document_id string    `gorm:"primarykey;default:gen_random_uuid()" json:"document_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image_url   string    `json:"image_url"`
	Status      string    `gorm:"check:status IN ('active', 'deleted')" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type DocumentRequest struct {
	Title       string                `form:"title" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required" swaggerignore:"true"`
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
