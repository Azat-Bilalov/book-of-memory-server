package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetActiveDocuments() ([]*ds.Document, error) {
	documents := make([]*ds.Document, 0)

	err := r.db.Find(&documents, "status = ?", ds.DOCUMENT_STATUS_ACTIVE).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) GetActiveDocumentByID(id int) (*ds.Document, error) {
	document := &ds.Document{}

	err := r.db.First(document, "document_id = ? AND status = ?", id, ds.DOCUMENT_STATUS_ACTIVE).Error
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (r *Repository) CreateDocument(document ds.Document) error {
	return r.db.Create(document).Error
}

func (r *Repository) UpdateDocument(document ds.Document) error {
	return r.db.Save(document).Error
}

func (r *Repository) DeleteDocument(id int) error {
	return r.db.Exec("UPDATE documents SET status = ? WHERE document_id = ?", ds.DOCUMENT_STATUS_DELETED, id).Error
}
