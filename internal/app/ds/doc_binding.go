package ds

type DocBinding struct {
	BindingID  uint `gorm:"primaryKey"`
	DocumentID uint `gorm:"primaryKey"`
}
