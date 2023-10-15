package ds

type DocBinding struct {
	Binding_id  uint `gorm:"primaryKey"`
	Document_id uint `gorm:"primaryKey"`
}
