package ds

type DocBinding struct {
	Binding_id  string `gorm:"primaryKey"`
	Document_id string `gorm:"primaryKey"`
}
