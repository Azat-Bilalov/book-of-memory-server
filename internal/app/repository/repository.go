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

func (r *Repository) GetDocuments() ([]*ds.Document, error) {
	documents := make([]*ds.Document, 0)

	err := r.db.Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *Repository) GetDocumentByID(id int) (*ds.Document, error) {
	document := &ds.Document{}

	err := r.db.First(document, "document_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (r *Repository) CreateDocument(document ds.Document) error {
	return r.db.Create(document).Error
}
