package ds

const (
	DOCUMENT_STATUS_ACTIVE  = "active"
	DOCUMENT_STATUS_DELETED = "deleted"
)

type Document struct {
	Document_id uint `gorm:"primarykey"`
	Title       string
	Description string
	Image_url   string
	Status      string
	Created_at  string
}
