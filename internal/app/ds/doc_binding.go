package ds

type DocBinding struct {
	Binding_id  string  `gorm:"primaryKey"`
	Document_id string  `gorm:"primaryKey"`
	Info        *string `gorm:"type:text"`
	File_url    *string
}

type DocBindingRequest struct {
	Info     *string `json:"info"`
	File_url *string `json:"file_url"`
}
