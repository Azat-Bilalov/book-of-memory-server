package repository

import (
	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
	"gorm.io/gorm"
)

type InterfaceDocBindingRepository interface {
	Store(docBinding *ds.DocBinding) (*ds.DocBinding, error)
	Find(documentID string, bindingID string) (*ds.DocBinding, error)
	Update(docBinding *ds.DocBinding) error
	Delete(documentID string, bindingID string) error
}

type DocBindingRepository struct {
	db *gorm.DB
}

func NewDocBindingRepository(db *gorm.DB) InterfaceDocBindingRepository {
	return &DocBindingRepository{db}
}

func (r *DocBindingRepository) Store(docBinding *ds.DocBinding) (*ds.DocBinding, error) {
	if err := r.db.Create(docBinding).Error; err != nil {
		return nil, err
	}

	return docBinding, nil
}

func (r *DocBindingRepository) Find(documentID string, bindingID string) (*ds.DocBinding, error) {
	docBinding := &ds.DocBinding{}
	err := r.db.First(docBinding, "document_id = ? AND binding_id = ?", documentID, bindingID).Error
	if err != nil {
		return nil, err
	}
	return docBinding, nil
}

func (r *DocBindingRepository) Update(docBinding *ds.DocBinding) error {
	return r.db.Save(docBinding).Error
}

func (r *DocBindingRepository) Delete(documentID string, bindingID string) error {
	return r.db.Exec("DELETE FROM doc_bindings WHERE document_id = ? AND binding_id = ?", documentID, bindingID).Error
}
