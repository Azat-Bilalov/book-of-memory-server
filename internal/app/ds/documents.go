package ds

type Document struct {
	Document_id uint `gorm:"primarykey"`
	Title       string
	Description string
	Image_url   string
	Status      string
	Created_at  string
}
