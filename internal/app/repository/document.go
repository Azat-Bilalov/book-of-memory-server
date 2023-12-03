package repository

import (
	"log"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"gorm.io/gorm"
)

type InterfaceDocumentRepository interface {
	Store(document *ds.Document) (*ds.Document, error)
	FindAll(status string, title string) ([]*ds.Document, error)
	FindByUUID(uuid string) (*ds.Document, error)
	CountWithImageUrl(imageUrl string) (int64, error)
	UpdateByUUID(document *ds.Document) (*ds.Document, error)
}

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) InterfaceDocumentRepository {
	return &DocumentRepository{db}
}

// Store сохраняет документ в базе данных
func (r *DocumentRepository) Store(document *ds.Document) (*ds.Document, error) {
	if err := r.db.Create(document).Error; err != nil {
		return nil, err
	}

	return document, nil
}

// FindAll возвращает все документы из базы данных
func (r *DocumentRepository) FindAll(status string, title string) ([]*ds.Document, error) {
	documents := make([]*ds.Document, 0)
	log.Println("FindAll")
	query := r.db.Table("documents").Where("status = ?", status).Where("lower(title) LIKE ?", "%"+title+"%")
	log.Println(query)
	if err := query.Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

// FindByUUID возвращает документ из базы данных по UUID
func (r *DocumentRepository) FindByUUID(uuid string) (*ds.Document, error) {
	document := &ds.Document{}

	err := r.db.First(document, "Document_id = ?", uuid).Error
	if err != nil {
		return nil, err
	}

	return document, nil
}

// CountWithImageUrl возвращает количество документов с такой же ссылкой на изображение
func (r *DocumentRepository) CountWithImageUrl(imageUrl string) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Document{}).Where("image_url = ?", imageUrl).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateByUUID обновляет документ в базе данных по UUID
func (r *DocumentRepository) UpdateByUUID(document *ds.Document) (*ds.Document, error) {
	if err := r.db.Save(document).Error; err != nil {
		return nil, err
	}

	return document, nil
}

// AddDocumentToBindingByUUID добавляет документ к привязке по UUID
// func (r *DocumentRepository) AddDocumentToBindingByUUID(uuid string, userID string, docBinding ds.DocBindingRequest) error {
// 	binding := &ds.Binding{}
// 	err := r.db.First(binding, "Binding_id = ?", uuid).Error
// 	if err != nil {
// 		return err
// 	}
// 	document := &ds.Document{}
// 	err = r.db.First(document, "Document_id = ?", docBinding.DocumentID).Error
// 	if err != nil {
// 		return err
// 	}
// 	if binding.UserID != userID {
// 		return gorm.ErrRecordNotFound
// 	}
// 	if binding.Status != ds.BINDING_STATUS_IN_PROGRESS {
// 		return gorm.ErrRecordNotFound
// 	}
// 	if err := r.db.Model(binding).Association("Documents").Append(document); err != nil {
// 		return err
// 	}
// 	return nil
// }
